package main

import (
	"fmt"
	"github.com/docker/libcontainer/user"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if os.Args[0] == "/init" {
		fmt.Fprintf(os.Stdout, "started dockersh persistent container\n")
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
	realUsername := tmplConfigVar(config.ContainerUsername, &configInterpolations)
	realHomedirTo := tmplConfigVar(config.MountHomeTo, &configInterpolations)
	realHomedirFrom := tmplConfigVar(config.MountHomeFrom, &configInterpolations)
	realImageName := tmplConfigVar(config.ImageName, &configInterpolations)
	realShell := tmplConfigVar(config.Shell, &configInterpolations)
	realUserCwd := tmplConfigVar(config.UserCwd, &configInterpolations)
	realContainerName := tmplConfigVar(config.ContainerName, &configInterpolations)

	_, err = dockerpid(realContainerName)
	if err != nil {
		_, err = dockerstart(realUsername, realHomedirFrom, realHomedirTo, realContainerName, realImageName, config.DockerSocket, config.MountHome, config.MountTmp, config.MountDockerSocket, config.Entrypoint, config.Cmd, config.DockerOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	_, _, groups, _, err := user.GetUserGroupSupplementaryHome(username, 65536, 65536, "/")
	err = nsenterexec(realContainerName, uid, gid, groups, realUserCwd, realShell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting shell in new container: %v\n", err)
		return 1
	}
	return 0
}
