package main

import (
	"strconv"
	"strings"
	"os"
	"os/exec"
)

func nsenterdetect() (found bool, err error) {
	cmd := exec.Command("boot2docker",  "ssh",  "[ -f /var/lib/boot2docker/nsenter ]")
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

func nsenterexec(pid int) (err error) {
	cmd := exec.Command("boot2docker", "ssh", "-t", "sudo", "/var/lib/boot2docker/nsenter",
			    "--target", strconv.Itoa(pid), "--mount", "--uts", "--ipc", "--net", "--pid",
			    "--", "/bin/ash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}
