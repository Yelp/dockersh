package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	found, err := nsenterdetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cound not detect if nsenter was installed: %s\n", err)
		return 1
	}
	if !found {
		fmt.Fprintf(os.Stderr, "nsenter is not installed\n")
		fmt.Fprintf(os.Stderr, "run boot2docker ssh 'docker run --rm -v /var/lib/boot2docker/:/target bobtfish/nsenter'\n")
		return 1
	}
	/* Woo! We found nsenter, now to move onto more interesting things */
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v", err)
		return 1
	}
	if user.HomeDir == "" {
		fmt.Fprintf(os.Stderr, "didn't get a home directory")
		return 1
	}
	if user.Username == "" {
		fmt.Fprintf(os.Stderr, "didn't get a username")
		return 1
	}

	containerName := fmt.Sprintf("%s_dockersh", user.Username)

	pid, err := dockerpid(containerName)
	if err != nil {
		pid, err = dockerstart(user.Username, user.HomeDir, containerName, "busybox")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	uid, err := strconv.Atoi(user.Uid)
	gid, err := strconv.Atoi(user.Gid)
	nsenterexec(pid, uid, gid, user.HomeDir, "/bin/sh")
	return 0
}
