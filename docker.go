package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func dockerpid(name string) (pid int, err error, out string) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Pid}}", name)
	output, err := cmd.Output()
	if err != nil {
		return -1, err, string(output)
	}

	pid, e := strconv.Atoi(strings.TrimSpace(string(output)))

	if e != nil {
		return -1, e, string(output)
	}
	if pid == 0 {
		return -1, errors.New("Invalid PID"), string(output)
	}
	return pid, nil, string(output)
}

func dockerstart(username string, homedir string, name string, container string) (pid int, err error, out string) {
	cmd := exec.Command("docker", "rm", "--name", name)
	err = cmd.Run()

	// docker run -t -i -u $DESIRED_USER --hostname="$MYHOSTNAME" --name="$DOCKER_NAME" -v $HOMEDIR:$HOMEDIR:rw -v /etc/passwd:/etc/passwd:ro -v /etc/group:/etc/group:ro -d "$DOCKER_CONTAINER"
	cmd = exec.Command("docker", "run", "-t", "-i", "-u", username, "-v", fmt.Sprintf("%s:%s:rw", homedir, homedir), "-v", "/etc/passwd:/etc/passwd:ro", "-v", "/etc/group:/etc/group:ro", "--name", name, "-d", container)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		return -1, err, output.String()
	}
	return dockerpid(name)
}
