package main

import (
	"fmt"
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
	config, err := loadAllConfig(username, homedir)
	configInterpolations := configInterpolation{homedir, username}
	realUsername := tmplConfigVar(config.ContainerUsername, &configInterpolations)
	realHomedirTo := tmplConfigVar(config.MountHomeTo, &configInterpolations)
	realHomedirFrom := tmplConfigVar(config.MountHomeFrom, &configInterpolations)
	realImageName := tmplConfigVar(config.ImageName, &configInterpolations)
	realShell := tmplConfigVar(config.Shell, &configInterpolations)
	containerName := fmt.Sprintf("%s_dockersh", realUsername)

	pid, err := dockerpid(containerName)
	if err != nil {
		pid, err = dockerstart(realUsername, realHomedirFrom, realHomedirTo, containerName, realImageName, config.DockerSocket, config.MountHome, config.MountTmp, config.MountDockerSocket, config.Entrypoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	err = nsenterexec(pid, uid, gid, realHomedirTo, realShell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return 1
	}
	return 0
}
