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
	MountHome           bool `json:"mount_home"`
	MountTmp            bool `json:"mount_tmp"`
}

type configInterpolation struct {
	Home string
	User string
}

var defaultConfig = Configuration{ImageName: "ubuntu", MountHomeTo: "%h", ContainerUsername: "%u", Shell: "%s", MountHome: true, MountTmp: true, BlacklistUserConfig: []string{}, BlacklistSetup: false}

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
	localConfigFile.Close()
	return (loadConfigFromString(bytes, config))
}

func loadConfigFromString(bytes []byte, config *Configuration) (err error) {
	var localConfig map[string]interface{}
	err = json.Unmarshal(bytes, &localConfig)
	if err != nil {
		return
	}

	for k, v := range localConfig {
		data, ok := v.(string)
		if !ok {
			return errors.New("parse")
		}
		configAllowed := true
		for _, element := range config.BlacklistUserConfig {
			if k == element {
				configAllowed = false
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
			case "shell":
				config.Shell = data
			case "blacklist_user_config":
				if !config.BlacklistSetup {
					for _, st := range strings.Split(data, ",") {
						config.BlacklistUserConfig = append(config.BlacklistUserConfig, st)
					}
					config.BlacklistSetup = true
				}
			}
		}
	}
	return nil
}
