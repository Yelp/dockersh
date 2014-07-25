package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

type configuration struct {
	ImageName           string   `json:"image_name"`
	MountHomeTo         string   `json:"mount_home_to"`
	ContainerUsername   string   `json:"container_username"`
	Shell               string   `json:"shell"`
	BlacklistUserConfig []string `json:"blacklist_user_config"`
}

func loadConfig() (config *configuration, err error) {
	config = &configuration{
		ImageName:         "busybox",
		MountHomeTo:       "{{Home}}",
		ContainerUsername: "{{User}}",
		Shell:             "/bin/ash",
	}
	localConfigFile, err := os.Open("dockersh.json")
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
		switch k {
		case "image_name":
			if localImageName, ok := v.(string); !ok {
				return nil, errors.New("parse")
			} else {
				config.ImageName = localImageName
			}
		case "mount_home_to":
			if localMountHomeTo, ok := v.(string); !ok {
				return nil, errors.New("parse")
			} else {
				config.MountHomeTo = localMountHomeTo
			}
		case "container_username":
			if localContainerUsername, ok := v.(string); !ok {
				return nil, errors.New("parse")
			} else {
				config.ContainerUsername = localContainerUsername
			}
		case "shell":
			if localShell, ok := v.(string); !ok {
				return nil, errors.New("parse")
			} else {
				config.Shell = localShell
			}
		}
	}
	return config, nil
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	found, err := nsenterdetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cound not detect if nsenter was installed: %s\n", err)
		return 1
	}
	if !found {
		fmt.Fprintf(os.Stderr, "nsenter is not installed\n")
		fmt.Fprintf(os.Stderr, "run boot2docker ssh 'docker run --rm -v /var/lib/boot2docker/:/target bobtfish/nsenter'\n")
		return 1
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	/* Woo! We found nsenter (if needed for this OS), now to move onto more interesting things */
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v", err)
		return 1
	}
	if user.HomeDir == "" {
		fmt.Fprintf(os.Stderr, "didn't get a home directory")
		return 1
	}
	if user.Username == "" {
		fmt.Fprintf(os.Stderr, "didn't get a username")
		return 1
	}

	containerName := fmt.Sprintf("%s_dockersh", user.Username)

	pid, err := dockerpid(containerName)
	if err != nil {
		pid, err = dockerstart(user.Username, user.HomeDir, containerName, config.ImageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	uid, err := strconv.Atoi(user.Uid)
	gid, err := strconv.Atoi(user.Gid)
	nsenterexec(pid, uid, gid, user.HomeDir, "/bin/sh")
	return 0
}
