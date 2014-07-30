package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type configuration struct {
	ImageName           string   `json:"image_name"`
	MountHomeTo         string   `json:"mount_home_to"`
	ContainerUsername   string   `json:"container_username"`
	Shell               string   `json:"shell"`
	BlacklistUserConfig []string `json:"blacklist_user_config"`
}

type configInterpolation struct {
	Home string
	User string
}

func loadConfig(filename string) (config *configuration, found bool, err error) {
	found = true
	config = &configuration{
		ImageName:         "ubuntu",
		MountHomeTo:       "%h",
		ContainerUsername: "%u",
		Shell:             "%s",
	}
	localConfigFile, err := os.Open(filename)
	if err != nil {
		found = false
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
			return nil, found, errors.New("parse")
		}
		switch k {
		case "image_name":
			config.ImageName = data
		case "mount_home_to":
			config.MountHomeTo = data
		case "container_username":
			config.ContainerUsername = data
		case "shell":
			config.Shell = data
		}
	}
	return config, found, nil
}

func main() {
	os.Exit(realMain())
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
	config, found, err := loadConfig(fmt.Sprintf("%s/.dockersh.json", homedir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	if !found {
		config, found, err = loadConfig("/etc/dockersh.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not load config: %v", err)
			return 1
		}
	}
	configInterpolations := configInterpolation{homedir, username}

	realUsername := tmplConfigVar(config.ContainerUsername, &configInterpolations)
	realHomedir := tmplConfigVar(config.MountHomeTo, &configInterpolations)
	realImageName := tmplConfigVar(config.ImageName, &configInterpolations)
	realShell := tmplConfigVar(config.Shell, &configInterpolations)
	containerName := fmt.Sprintf("%s_dockersh", realUsername)

	pid, err := dockerpid(containerName)
	if err != nil {
		pid, err = dockerstart(realUsername, realHomedir, containerName, realImageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	nsenterexec(pid, uid, gid, realHomedir, realShell)
	return 0
}
