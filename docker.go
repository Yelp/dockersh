package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func dockerpid(name string) (pid int, err error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Pid}}", name)
	output, err := cmd.Output()
	if err != nil {
		return -1, errors.New(err.Error() + ":\n" + string(output))
	}

	pid, err = strconv.Atoi(strings.TrimSpace(string(output)))

	if err != nil {
		return -1, errors.New(err.Error() + ":\n" + string(output))
	}
	if pid == 0 {
		return -1, errors.New("Invalid PID")
	}
	return pid, nil
}

func dockerstart(username string, homedir string, name string, container string) (pid int, err error) {
	cmd := exec.Command("docker", "rm", name)
	err = cmd.Run()

	this_binary, _ := filepath.Abs(os.Args[0])
	// FIXME - Binding /tmp to host, can we get ssh working a better way?
	cmd = exec.Command("docker", "run", "-d", "-u", username, "-v", fmt.Sprintf("%s:%s:rw", homedir, homedir),
		"-v", "/tmp:/tmp", "-v", "/etc/passwd:/etc/passwd:ro", "-v", "/etc/group:/etc/group:ro",
		"-v", this_binary+":/sbin/init", "--name", name,
		"--entrypoint", "/sbin/init", container, "--")

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		return -1, errors.New(err.Error() + ":\n" + output.String())
	}
	return dockerpid(name)
}
