package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func nsenterdetect() (found bool, err error) {
	cmd := exec.Command("boot2docker", "ssh", "[ -f /var/lib/boot2docker/nsenter ]")
	err = cmd.Run()
	if err == nil {
		return true, nil
	}
	/* TODO: Figure out how to get the actual error code from here */
	if e, ok := err.(*exec.ExitError); ok && strings.HasSuffix(e.String(), "1") {
		return false, nil
	}
	return false, err
}

func nsenterexec(pid int, uid int, gid int, wd string, shell string) (err error) {
	// sudo nsenter --target "$PID" --mount --uts --ipc --net --pid --setuid $DESIRED_UID --setgid $DESIRED_GID --wd=$HOMEDIR -- "$REAL_SHELL"
	cmd := exec.Command("boot2docker", "ssh", "-t", "sudo", "/var/lib/boot2docker/nsenter",
		"--target", strconv.Itoa(pid), "--mount", "--uts", "--ipc", "--net", "--pid",
		"--setuid", strconv.Itoa(uid), "--setgid", strconv.Itoa(gid), fmt.Sprintf("--wd=%s", wd),
		"--", shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}
