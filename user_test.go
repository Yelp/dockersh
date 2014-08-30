package main

import "testing"
import "os/user"

func Test_Add2Ints_1(t *testing.T) {
	mockuser := &user.User{Username: "vagrant", HomeDir: "/home/vagrant", Uid: "1000", Gid: "1000"}
	username, homedir, uid, gid, err := getUser(mockuser)
	if err != nil {
		t.Error("Got error from getUser " + err.Error())
	}
	if username == "vagrant" {
		t.Log("username passed.")
	} else {
		t.Errorf("Username failed: %s", username)
	}
	if homedir == "/home/vagrant" {
		t.Log("homedir passed.")
	} else {
		t.Errorf("homedir failed: %s", homedir)
	}
	if uid == 1000 {
		t.Log("uid passed.")
	} else {
		t.Errorf("uid failed: %i", uid)
	}
	if gid == 1000 {
		t.Log("git passed.")
	} else {
		t.Errorf("gid failed: %i", gid)
	}
}
