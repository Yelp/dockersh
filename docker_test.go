package main

import (
	"testing"
)

func Test_dockerPid_1(t *testing.T) {
	pid, err := dockerpid("testcontainer")
	if err != nil {
		t.Errorf("Error from dockerpid: %v", err)
	}
	if pid != 666 {
		t.Errorf("PID was %i expected 666", pid)
	}
}

func Test_dockerSha_1(t *testing.T) {
	sha, err := dockersha("testcontainer")
	if err != nil {
		t.Errorf("Error from dockersha: %v", err)
	}
	if sha != "666" {
		t.Errorf("SHA was %s expected 666", sha)
	}
}

func Test_dockerStart(t *testing.T) {
	c := Configuration{ContainerName: "somecontainer", ImageName: "busybox", MountHome: true, MountHomeFrom: "/home/fred", MountHomeTo: "/home/fred", Entrypoint: "internal", DockerSocket: "dockersock", Cmd: []string{"foo"}, DockerOpt: []string{"bar"}}
	pid, err := dockerstart(c)
	if err != nil {
		t.Errorf("Error from dockerstart: %v", err)
	}
	if pid != 666 {
		t.Errorf("PID was %i expected 666", pid)
	}
}

func Test_validatePortforwardString_1(t *testing.T) {
	err := validatePortforwardString("1:2")
	if err != nil {
		t.Errorf("Error on 1:2")
	}
}

