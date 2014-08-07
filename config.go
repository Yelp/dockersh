package main

import (
	"code.google.com/p/gcfg"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	ImageName                   string
	EnableUserImageName         bool
	MountHomeFrom               string
	EnableUserMountHomeFrom     bool
	MountHomeTo                 string
	EnableUserMountHomeTo       bool
	ContainerUsername           string
	EnableUserContainerUsername bool
	Shell                       string
	EnableUserShell             bool
	EnableUserConfig            bool
	MountHome                   bool
	EnableUserMountHome         bool
	MountTmp                    bool
	EnableUserMountTmp          bool
	MountDockerSocket           bool
	EnableUserMountDockerSocket bool
	DockerSocket                string
	EnableUserDockerSocket      bool
	Entrypoint                  string
	EnableUserEntrypoint        bool
}

func (c Configuration) Dump() string {
	return fmt.Sprintf("ImageName %s MountHomeTo %s ContainerUsername %s Shell %s DockerSocket %s", c.ImageName, c.MountHomeTo, c.ContainerUsername, c.Shell, c.DockerSocket)
}

type configInterpolation struct {
	Home string
	User string
}

var defaultConfig = Configuration{
	ImageName:         "busybox",
	MountHomeFrom:     "%h",
	MountHomeTo:       "%h",
	ContainerUsername: "%u",
	Shell:             "/bin/ash",
	DockerSocket:      "/var/run/docker.sock",
	Entrypoint:        "internal",
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
	if (!blacklist || old.EnableUserShell) && new.Shell != "" {
		old.Shell = new.Shell
	}
	if (!blacklist || old.EnableUserContainerUsername) && new.ContainerUsername != "" {
		old.ContainerUsername = new.ContainerUsername
	}
	if (!blacklist || old.EnableUserImageName) && new.ImageName != "" {
		old.ImageName = new.ImageName
	}
	if (!blacklist || old.EnableUserMountHomeTo) && new.MountHomeTo != "" {
		old.MountHomeTo = new.MountHomeTo
	}
	if (!blacklist || old.EnableUserMountHomeFrom) && new.MountHomeFrom != "" {
		old.MountHomeFrom = new.MountHomeFrom
	}
	if (!blacklist || old.EnableUserDockerSocket) && new.DockerSocket != "" {
		old.DockerSocket = new.DockerSocket
	}
	if (!blacklist || old.EnableUserMountHome) && new.MountHome == true {
		old.MountHome = true
	}
	if (!blacklist || old.EnableUserMountTmp) && new.MountTmp == true {
		old.MountTmp = true
	}
	if (!blacklist || old.EnableUserMountDockerSocket) && new.MountDockerSocket == true {
		old.MountDockerSocket = true
	}
	if (!blacklist || old.EnableUserEntrypoint) && new.Entrypoint != "" {
		old.Entrypoint = new.Entrypoint
	}
	if !blacklist && new.EnableUserConfig == true {
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
