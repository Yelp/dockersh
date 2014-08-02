package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	fmt.Fprintf(os.Stdout, "starting dockersh root process\n")
	if os.Args[0] == "/sbin/init" {
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
	var config = defaultConfig
	err = loadConfig("/etc/dockersh.json", &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	err = loadConfig(fmt.Sprintf("%s/.dockersh.json", homedir), &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load config: %v", err)
		return 1
	}
	configInterpolations := configInterpolation{homedir, username}
	realUsername := tmplConfigVar(config.ContainerUsername, &configInterpolations)
	realHomedir := tmplConfigVar(config.MountHomeTo, &configInterpolations)
	realImageName := tmplConfigVar(config.ImageName, &configInterpolations)
	realShell := tmplConfigVar(config.Shell, &configInterpolations)
	containerName := fmt.Sprintf("%s_dockersh", realUsername)

	pid, err := dockerpid(containerName)
	if err != nil {
		// bools are bindtmp, bindhome, last string is the init process
		pid, err = dockerstart(realUsername, realHomedir, containerName, realImageName, true, true, "internal")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	nsenterexec(pid, uid, gid, realHomedir, realShell)
	return 0
}
