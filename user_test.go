package main

import "testing"
import "os/user"
import "fmt"

func Test_Add2Ints_1(t *testing.T) {
	mockuser, err := user.Current() // Worst mock evar
	if err != nil {
		t.Error(fmt.Sprintf("could not get current user: %v", err))
	}
	username, homedir, uid, gid, err := getUser(mockuser)
	if username == "vagrant" {
  		t.Log("username passed.")
	} else {
		t.Error(fmt.Sprintf("Username failed: %s", username))
	}
	if homedir == "/home/vagrant" {
                t.Log("homedir passed.")
        } else {
                t.Error(fmt.Sprintf("homedir failed: %s", homedir))
	}
	if uid == 1000 {
		t.Log("uid passed.")
	} else {
		t.Error(fmt.Sprintf("uid failed: %i", uid))
	}
	if gid == 1000 {
		t.Log("git passed.")
	} else {
		t.Error(fmt.Sprintf("gid failed: %i", gid))
	}
}

