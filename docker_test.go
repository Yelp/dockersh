package main

import (
	"testing"
	"fmt"
)

func Test_dockerPid_1(t *testing.T) {
	pid, err := dockerpid("testcontainer")
	if err != nil {
		t.Error(fmt.Sprintf("Error from dockerpid: %v", err))
	}
	if pid != 666 {
		t.Error(fmt.Sprintf("PID was %i expected 666", pid))
	}
}

func Test_dockerSha_1(t *testing.T) {
	sha, err := dockersha("testcontainer")
	if err != nil {
		t.Error(fmt.Sprintf("Error from dockersha: %v", err))
	}
	if sha != "666" {
		t.Error(fmt.Sprintf("SHA was %s expected 666", sha))
	}
}

func Test_dockerStart(t *testing.T) {
	pid, err := dockerstart("someuser", "homedirfrom", "homedirto", "name", "container", "dockersock", true, true, true, "internal", []string{"foo"}, []string{"bar"})
	if err != nil {
                t.Error(fmt.Sprintf("Error from dockerstart: %v", err))
	}
	if pid != 666 {
                t.Error(fmt.Sprintf("PID was %i expected 666", pid))
        }
}


