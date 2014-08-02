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

func dockerstart(username string, homedir string, name string, container string, dockersock string, bindtmp bool, bindhome bool, binddocker bool, init string) (pid int, err error) {
	cmd := exec.Command("docker", "rm", name)
	err = cmd.Run()

	bind_self_as_init := false
	if init == "internal" {
		init = "/sbin/init"
		bind_self_as_init = true
	}
	this_binary := "/usr/local/bin/dockersh"
	if os.Getenv("SHELL") != "/usr/local/bin/dockersh" {
		this_binary, _ = filepath.Abs(os.Args[0])
	}
	// FIXME - Binding /tmp to host, can we get ssh working a better way?
	var cmdtxt = []string{"run", "-d", "-u", username,
		"-v", "/etc/passwd:/etc/passwd:ro", "-v", "/etc/group:/etc/group:ro"}

	if bindtmp {
		cmdtxt = append(cmdtxt, "-v", "/tmp:/tmp")
	}
	if bindhome {
		cmdtxt = append(cmdtxt, "-v", fmt.Sprintf("%s:%s:rw", homedir, homedir))
	}
	if bind_self_as_init {
		fmt.Fprintf(os.Stderr, "This binary is %s\n", this_binary)
		cmdtxt = append(cmdtxt, "-v", this_binary+":/sbin/init")
	}
	if binddocker {
		cmdtxt = append(cmdtxt, "-v", dockersock+":/var/run/docker.sock")
	}
	cmdtxt = append(cmdtxt, "--name", name, "--entrypoint", init, container, "--")
	cmd = exec.Command("docker", cmdtxt...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		return -1, errors.New(err.Error() + ":\n" + output.String())
	}
	return dockerpid(name)
}
