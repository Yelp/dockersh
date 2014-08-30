package main

import "testing"
import "os/user"

func Test_getCurrentUser_1(t *testing.T) {
    _, _, _, _, err := getCurrentUser()
    if err != nil {
        t.Error("Error from getCurrentUser")
    }
}

func Test_getUser_1(t *testing.T) {
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

func Test_getUser_2(t *testing.T) {
	mockuser := &user.User{Username: "", HomeDir: "/home/vagrant", Uid: "1000", Gid: "1000"}
	_, _, _, _, err := getUser(mockuser)
	if err == nil {
		t.Error("No error from getUser")
	}
}

func Test_getUser_3(t *testing.T) {
	mockuser := &user.User{Username: "Foo", HomeDir: "", Uid: "1000", Gid: "1000"}
	_, _, _, _, err := getUser(mockuser)
	if err == nil {
		t.Error("No error from getUser")
	}
}

