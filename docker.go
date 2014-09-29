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

func dockerVersionCheck() (err error) {
	versionString, err := getDockerVersionString()
	// Docker version 1.1.2, build d84a070
	versionStringParts := strings.Split(versionString, " ")
	versionParts := strings.Split(versionStringParts[2], ".")
	major, _ := strconv.Atoi(versionParts[0])
	minor, _ := strconv.Atoi(versionParts[1])
	if major > 1 {
		return nil
	}
	if minor >= 2 {
		return nil
	}
	return errors.New(fmt.Sprintf("Docker version '%s' lower than desired version '1.2.0'", versionStringParts[2]))
}

func getDockerVersionString() (string, error) {
	cmd := exec.Command("docker", "-v")
	o, err := cmd.Output()
	return string(o), err
}

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

func dockerstart(config Configuration) (pid int, err error) {
	cmd := exec.Command("docker", "rm", config.ContainerName)
	_ = cmd.Run()
	cmdtxt, err := dockercmdline(config)
	if err != nil {
		return -1, err
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
	return dockerpid(config.ContainerName)
}

func dockercmdline(config Configuration) ([]string, error) {
	var err error
	bindSelfAsInit := false
	init := config.Entrypoint
	if init == "internal" {
		init = "/init"
		bindSelfAsInit = true
	}
	thisBinary := "/usr/local/bin/dockersh"
	if os.Getenv("SHELL") != "/usr/local/bin/dockersh" {
		thisBinary, _ = filepath.Abs(os.Args[0])
	}
	var cmdtxt = []string{"run", "-d", "-u", config.ContainerUsername,
		"-v", "/etc/passwd:/etc/passwd:ro", "-v", "/etc/group:/etc/group:ro",
		"--cap-drop", "SETUID", "--cap-drop", "SETGID", "--cap-drop", "NET_RAW",
		"--cap-drop", "MKNOD"}
	if len(config.DockerOpt) > 0 {
		for _, element := range config.DockerOpt {
			cmdtxt = append(cmdtxt, element)
		}
	}
	if config.MountTmp {
		cmdtxt = append(cmdtxt, "-v", "/tmp:/tmp")
	}
	if config.MountHome {
		cmdtxt = append(cmdtxt, "-v", fmt.Sprintf("%s:%s:rw", config.MountHomeFrom, config.MountHomeTo))
	}
	if bindSelfAsInit {
		cmdtxt = append(cmdtxt, "-v", thisBinary+":/init")
	} else {
		if len(config.ReverseForward) > 0 {
			return []string{}, errors.New("Cannot configure ReverseForward with a custom init process")
		}
	}
	if config.MountDockerSocket {
		cmdtxt = append(cmdtxt, "-v", config.DockerSocket+":/var/run/docker.sock")
	}
	if len(config.ReverseForward) > 0 {
		cmdtxt, err = setupReverseForward(cmdtxt, config.ReverseForward)
		if err != nil {
			return []string{}, err
		}
	}
	cmdtxt = append(cmdtxt, "--name", config.ContainerName, "--entrypoint", init, config.ImageName)
	if len(config.Cmd) > 0 {
		for _, element := range config.Cmd {
			cmdtxt = append(cmdtxt, element)
		}
	} else {
		cmdtxt = append(cmdtxt, "")
	}

	return cmdtxt, nil
}

func validatePortforwardString(element string) error {
	parts := strings.Split(element, ":")
	if len(parts) != 2 {
		return errors.New("Number of parts must be 2")
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return(err)
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return(err)
	}
	return nil
}

func setupReverseForward(cmdtxt []string, reverseForward []string) ([]string, error) {
	for _, element := range reverseForward {
		err := validatePortforwardString(element)
		if err != nil {
			return cmdtxt, err
		}
	}
	cmdtxt = append(cmdtxt, "--env=DOCKERSH_PORTFORWARD="+strings.Join(reverseForward, ","))
	return cmdtxt, nil
}
