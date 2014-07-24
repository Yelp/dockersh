package main

import (
	"fmt"
	"os"
	"os/user"
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
	u, err2 := user.Current()
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Current: %v", err2)
	}
	if u.HomeDir == "" {
		fmt.Fprintf(os.Stderr, "didn't get a HomeDir")
	}
	if u.Username == "" {
		fmt.Fprintf(os.Stderr, "didn't get a username")
	}

	var container_name = fmt.Sprintf("%s_dockersh", u.Username)

	pid, err, out := dockerpid(container_name)
	if err != nil {
		pid, err, out = dockerstart(u.Username, u.HomeDir, container_name, "busybox")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cound not start container: %s: %s\n", err, out)
			return 1
		}
	}
	nsenterexec(pid)
	return 0
}
