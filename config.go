package main

import (
	"code.google.com/p/gcfg"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

func (c Configuration) Dump() string {
	return fmt.Sprintf("ImageName %s MountHomeTo %s ContainerUsername %s Shell %s DockerSocket %s", c.ImageName, c.MountHomeTo, c.ContainerUsername, c.Shell, c.DockerSocket)
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

func loadAllConfig(user string, homedir string) (config Configuration, err error) {
	globalconfig, err := loadConfig(loadableFile("/etc/dockersh"), user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return config, errors.New("could not load config")
	}
	localconfig, err := loadConfig(loadableFile(fmt.Sprintf("%s/.dockersh", homedir)), user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return config, errors.New("could not load config")
	}
	return mergeConfigs(mergeConfigs(defaultConfig, globalconfig, false), localconfig, true), nil
}

type loadableFile string

func (fn loadableFile) Getcontents() []byte {
	localConfigFile, err := os.Open(string(fn))
	if err != nil {
	}
	b, err := ioutil.ReadAll(localConfigFile)
	localConfigFile.Close()
	return b
}

func loadConfig(filename loadableFile, user string) (config Configuration, err error) {
	bytes := filename.Getcontents()
	if err != nil {
		return
	}
	return (loadConfigFromString(bytes, user))
}

func mergeConfigs(old Configuration, new Configuration, blacklist bool) (ret Configuration) {
	if blacklist && old.DisableUserConfig {
		return old
	}
	var m = make(map[string]bool)
	for _, element := range old.BlacklistUserConfig {
		m[element] = true
	}
	if (!blacklist || !m["shell"]) && new.Shell != "" {
		old.Shell = new.Shell
	}
	if (!blacklist || !m["containerusername"]) && new.ContainerUsername != "" {
		old.ContainerUsername = new.ContainerUsername
	}
	if (!blacklist || !m["imagename"]) && new.ImageName != "" {
		old.ImageName = new.ImageName
	}
	if (!blacklist || !m["mounthometo"]) && new.MountHomeTo != "" {
		old.MountHomeTo = new.MountHomeTo
	}
	if (!blacklist || !m["dockersocket"]) && new.DockerSocket != "" {
		old.DockerSocket = new.DockerSocket
	}
	return old
}

func loadConfigFromString(bytes []byte, user string) (config Configuration, err error) {
	inicfg := struct {
		Dockersh Configuration
		User     map[string]*Configuration
	}{}
	err = gcfg.ReadStringInto(&inicfg, string(bytes))
	if err != nil {
		return
	}
	if inicfg.User[user] == nil {
		return inicfg.Dockersh, nil
	}
	return mergeConfigs(inicfg.Dockersh, *inicfg.User[user], false), nil
}
