package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	found, err := nsenterdetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cound not detect if nsenter was installed: %s\n", err);
		return 1
	}
	if !found {
		fmt.Fprintf(os.Stderr, "nsenter is not installed\n");
		fmt.Fprintf(os.Stderr, "run boot2docker ssh 'docker run --rm -v /var/lib/boot2docker/:/target jpetazzo/nsenter'\n");
		return 1
	}
	/* Woo! We found nsenter, now to move onto more interesting things */
	pid, err := dockerpid("juliank_shell")
	if err != nil {
		pid, err = dockerstart("juliank_shell", "busybox")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cound not start container: %s\n", err);
			return 1
		}
	}
	nsenterexec(pid)
	return 0
}
