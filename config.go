package main

import (
	"code.google.com/p/gcfg"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	ImageName           string
	MountHomeTo         string
	ContainerUsername   string
	Shell               string
	BlacklistUserConfig []string
	BlacklistSetup      bool
	EnableUserConfig    bool
	MountHome           bool
	MountTmp            bool
	MountDockerSocket   bool
	DockerSocket        string
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
	BlacklistUserConfig: []string{"imagename", "shell", "containerusername", "mounthometo", "mounttmp", "mountdockersocket"},
	BlacklistSetup:      false,
	EnableUserConfig:    false,
	MountDockerSocket:   false,
	DockerSocket:        "/var/run/docker.sock",
}

func loadAllConfig(user string, homedir string) (config Configuration, err error) {
	globalconfig, err := loadConfig(loadableFile("/etc/dockersh"), user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return config, errors.New("could not load config")
	}
	if globalconfig.EnableUserConfig == true {
		localconfig, err := loadConfig(loadableFile(fmt.Sprintf("%s/.dockersh", homedir)), user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not load config: %v", err)
			return config, errors.New("could not load config")
		}
		return mergeConfigs(mergeConfigs(defaultConfig, globalconfig, false), localconfig, true), nil
	} else {
		return mergeConfigs(defaultConfig, globalconfig, false), nil
	}

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
	if new.EnableUserConfig == true {
		old.EnableUserConfig = true
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
