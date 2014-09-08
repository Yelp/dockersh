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

func dockersha(name string) (sha string, err error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.Id}}", name)
	output, err := cmd.Output()
	if err != nil {
		return sha, errors.New(err.Error() + ":\n" + string(output))
	}
	sha = strings.TrimSpace(string(output))
	if sha == "" {
		return "", errors.New("Invalid SHA")
	}
	return sha, nil
}

func dockerstart(username string, homedirfrom string, homedirto string, name string, container string, dockersock string, bindhome bool, bindtmp bool, binddocker bool, init string, cmdargs []string, dockeropts []string) (pid int, err error) {
	cmd := exec.Command("docker", "rm", name)
	err = cmd.Run()

	bindSelfAsInit := false
	if init == "internal" {
		init = "/init"
		bindSelfAsInit = true
	}
	thisBinary := "/usr/local/bin/dockersh"
	if os.Getenv("SHELL") != "/usr/local/bin/dockersh" {
		thisBinary, _ = filepath.Abs(os.Args[0])
	}
	var cmdtxt = []string{"run", "-d", "-u", username,
		"-v", "/etc/passwd:/etc/passwd:ro", "-v", "/etc/group:/etc/group:ro",
		"--cap-drop", "SETUID", "--cap-drop", "SETGID", "--cap-drop", "NET_RAW",
		"--cap-drop", "MKNOD"}
	if len(dockeropts) > 0 {
		for _, element := range dockeropts {
			cmdtxt = append(cmdtxt, element)
		}
	}
	if bindtmp {
		cmdtxt = append(cmdtxt, "-v", "/tmp:/tmp")
	}
	if bindhome {
		cmdtxt = append(cmdtxt, "-v", fmt.Sprintf("%s:%s:rw", homedirfrom, homedirto))
	}
	if bindSelfAsInit {
		cmdtxt = append(cmdtxt, "-v", thisBinary+":/init")
	}
	if binddocker {
		cmdtxt = append(cmdtxt, "-v", dockersock+":/var/run/docker.sock")
	}
	cmdtxt = append(cmdtxt, "--name", name, "--entrypoint", init, container)
	if len(cmdargs) > 0 {
		for _, element := range cmdargs {
			cmdtxt = append(cmdtxt, element)
		}
	} else {
		cmdtxt = append(cmdtxt, "")
	}
	//fmt.Fprintf(os.Stderr, "docker %s\n", strings.Join(cmdtxt, " "))
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
