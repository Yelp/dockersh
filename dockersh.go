package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	_, err := nsenterdetect()
	if err != nil {
		return 1
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	/* Woo! We found nsenter, now to move onto more interesting things */
	username, homedir, uid, gid, err := getCurrentUser()

	containerName := fmt.Sprintf("%s_dockersh", username)

	pid, err := dockerpid(containerName)
	if err != nil {
		pid, err = dockerstart(username, homedir, containerName, config.ImageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	nsenterexec(pid, uid, gid, homedir, "/bin/ash")
	return 0
}
