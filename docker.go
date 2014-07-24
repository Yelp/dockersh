package main

import (
	"os/exec"
	"strconv"
	"strings"
	"errors"
)

func dockerpid(name string) (pid int, err error) {
	cmd := exec.Command("docker",  "inspect",  "--format", "{{.State.Pid}}", name)
	output, err := cmd.Output()
	if err != nil {
		return -1, err
	}

	pid, e := strconv.Atoi(strings.TrimSpace(string(output)))

	if e != nil {
		return -1, e
	}
	if pid == 0 {
		return -1, errors.New("Invalid PID")
	}
	return pid, nil
}


func dockerstart(name string, container string) (pid int, err error) {
	cmd := exec.Command("docker",  "run",  "-t", "-i", "--name", name, "-d", container)
	err = cmd.Run()
	if err != nil {
		return -1, err
	}
	return dockerpid(name)
}
