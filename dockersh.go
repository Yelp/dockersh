package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/libcontainer/user"
)

func main() {
	if os.Args[0] == "/init" {
		os.Exit(initMain())
	} else {
		os.Exit(realMain())
	}
}

func tmplConfigVar(template string, v *configInterpolation) string {
	shell := "/bin/bash"
	r := strings.NewReplacer("%h", v.Home, "%u", v.User, "%s", shell) // Arguments are old, new ...
	return r.Replace(template)
}

func getInterpolatedConfig(config *Configuration, configInterpolations configInterpolation) error {
	config.ContainerUsername = tmplConfigVar(config.ContainerUsername, &configInterpolations)
	config.MountHomeTo = tmplConfigVar(config.MountHomeTo, &configInterpolations)
	config.MountHomeFrom = tmplConfigVar(config.MountHomeFrom, &configInterpolations)
	config.ImageName = tmplConfigVar(config.ImageName, &configInterpolations)
	config.Shell = tmplConfigVar(config.Shell, &configInterpolations)
	config.UserCwd = tmplConfigVar(config.UserCwd, &configInterpolations)
	config.ContainerName = tmplConfigVar(config.ContainerName, &configInterpolations)
	return nil
}

func initMain() int {
	fmt.Fprintf(os.Stdout, "started dockersh persistent container\n")
	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc
	return 0
}

func realMain() int {
	username, homedir, uid, gid, err := getCurrentUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v", err)
		return 1
	}
	config, err := loadAllConfig(username, homedir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config: %v\n", err)
		return 1
	}
	configInterpolations := configInterpolation{homedir, username}
	err = getInterpolatedConfig(&config, configInterpolations)
	if err != nil {
		panic(fmt.Sprintf("Cannot interpolate config: %v", err))
	}

	_, err = dockerpid(config.ContainerName)
	if err != nil {
		_, err = dockerstart(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	_, _, groups, _, err := user.GetUserGroupSupplementaryHome(username, 65536, 65536, "/")
	err = nsenterexec(config.ContainerName, uid, gid, groups, config.UserCwd, config.Shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting shell in new container: %v\n", err)
		return 1
	}
	return 0
}
