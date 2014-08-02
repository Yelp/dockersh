package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	ImageName           string   `json:"image_name"`
	MountHomeTo         string   `json:"mount_home_to"`
	ContainerUsername   string   `json:"container_username"`
	Shell               string   `json:"shell"`
	BlacklistUserConfig []string `json:"blacklist_user_config"`
	BlacklistSetup      bool
	DisableUserConfig   bool   `json:"disable_user_config"`
	MountHome           bool   `json:"mount_home"`
	MountTmp            bool   `json:"mount_tmp"`
	MountDockerSocket   bool   `json:"mount_docker_socket"`
	DockerSocket        string `json:"docker_socket"`
}

type configInterpolation struct {
	Home string
	User string
}

var defaultConfig = Configuration{
	ImageName:           "busybox",
	MountHomeTo:         "%h",
	ContainerUsername:   "%u",
	Shell:               "/bin/ash",
	MountHome:           true,
	MountTmp:            true,
	BlacklistUserConfig: []string{"image_name", "shell", "container_username", "mount_home_to", "mount_tmp", "mount_docker_socket"},
	BlacklistSetup:      false,
	DisableUserConfig:   false,
	MountDockerSocket:   false,
	DockerSocket:        "/var/run/docker.sock",
}

func loadConfig(filename string, config *Configuration, limit bool) (err error) {
	localConfigFile, err := os.Open(filename)
	if err != nil {
		err = nil
		return
	}
	bytes, err := ioutil.ReadAll(localConfigFile)
	if err != nil {
		return
	}
	localConfigFile.Close()
	return (loadConfigFromString(bytes, config, limit))
}

func loadConfigFromString(bytes []byte, config *Configuration, limit bool) (err error) {
	var localConfig map[string]interface{}
	err = json.Unmarshal(bytes, &localConfig)
	if err != nil {
		return
	}
	if config.DisableUserConfig != false {
		return nil
	}
	for k, v := range localConfig {
		data, ok := v.(string)
		if !ok {
			return errors.New("parse")
		}
		configAllowed := true
		if limit {
			for _, element := range config.BlacklistUserConfig {
				if k == element {
					configAllowed = false
				}
			}
		}
		if configAllowed {
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
			case "disable_user_config":
				if data == "true" {
					config.DisableUserConfig = true
				}
			case "shell":
				config.Shell = data
			case "blacklist_user_config":
				if !config.BlacklistSetup {
					for _, st := range strings.Split(data, ",") {
						config.BlacklistUserConfig = append(config.BlacklistUserConfig, st)
					}
					config.BlacklistSetup = true
				}
			case "docker_socket":
				config.DockerSocket = data
			case "mount_docker_socket":
				if data == "true" {
					config.MountDockerSocket = true
				} else {
					config.MountDockerSocket = false
				}
			}
		}
	}
	return nil
}
