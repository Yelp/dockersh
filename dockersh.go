package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Configuration struct {
	ImageName           string   `json:"image_name"`
	MountHomeTo         string   `json:"mount_home_to"`
	ContainerUsername   string   `json:"container_username"`
	Shell               string   `json:"shell"`
	BlacklistUserConfig []string `json:"blacklist_user_config"`
	MountHome           bool     `json:"mount_home"`
	MountTmp            bool     `json:"mount_tmp"`
}

type configInterpolation struct {
	Home string
	User string
}

var defaultConfig = Configuration{ImageName: "ubuntu", MountHomeTo: "%h", ContainerUsername: "%u", Shell: "%s", MountHome: true, MountTmp: true}

func loadConfig(filename string, config *Configuration) (err error) {
	localConfigFile, err := os.Open(filename)
	if err != nil {
		err = nil
		return
	}
	bytes, err := ioutil.ReadAll(localConfigFile)
	if err != nil {
		return
	}
	var localConfig map[string]interface{}
	err = json.Unmarshal(bytes, &localConfig)
	if err != nil {
		return
	}
	localConfigFile.Close()

	for k, v := range localConfig {
		data, ok := v.(string)
		if !ok {
			return errors.New("parse")
		}
		switch k {
		case "image_name":
			config.ImageName = data
		case "mount_home_to":
			config.MountHomeTo = data
		case "container_username":
			config.ContainerUsername = data
		case "mount_tmp":
			if data == "true" {
				config.MountTmp = true
			} else {
				config.MountTmp = false
			}
		case "mount_home":
			if data == "true" {
				config.MountHome = true
			} else {
				config.MountHome = false
			}
		case "shell":
			config.Shell = data
		}
	}
	return nil
}

func main() {
	fmt.Fprintf(os.Stdout, "starting dockersh root process\n")
	if os.Args[0] == "/sbin/init" {
		// Wait for terminating signal
		sc := make(chan os.Signal, 2)
		signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
		<-sc
		os.Exit(0)
	} else {
		os.Exit(realMain())
	}
}

func tmplConfigVar(template string, v *configInterpolation) string {
	shell := "/bin/bash"
	return strings.Replace(strings.Replace(strings.Replace(template, "%h", v.Home, -1), "%u", v.User, -1), "%s", shell, -1)
}

func realMain() int {
	_, err := nsenterdetect()
	if err != nil {
		return 1
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	/* Woo! We found nsenter, now to move onto more interesting things */
	username, homedir, uid, gid, err := getCurrentUser()
	var config = defaultConfig
	err = loadConfig("/etc/dockersh.json", &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	err = loadConfig(fmt.Sprintf("%s/.dockersh.json", homedir), &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	configInterpolations := configInterpolation{homedir, username}
	realUsername := tmplConfigVar(config.ContainerUsername, &configInterpolations)
	realHomedir := tmplConfigVar(config.MountHomeTo, &configInterpolations)
	realImageName := tmplConfigVar(config.ImageName, &configInterpolations)
	realShell := tmplConfigVar(config.Shell, &configInterpolations)
	containerName := fmt.Sprintf("%s_dockersh", realUsername)

	pid, err := dockerpid(containerName)
	if err != nil {
		// bools are bindtmp, bindhome, last string is the init process
		pid, err = dockerstart(realUsername, realHomedir, containerName, realImageName, true, true, "internal")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	nsenterexec(pid, uid, gid, realHomedir, realShell)
	return 0
}
