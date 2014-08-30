package main

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
)

func getCurrentUser() (username string, homedir string, uid int, gid int, err error) {
	user, err := user.Current()
	if err != nil {
		return "", "", 0, 0, errors.New(fmt.Sprintf("could not get current user: %v", err))
	}
	return getUser(user)
}

func getUser(user *user.User) (username string, homedir string, uid int, gid int, err error) {
	if user.HomeDir == "" {
		return "", "", 0, 0, errors.New("didn't get a home directory")
	}
	if user.Username == "" {
		return "", "", 0, 0, errors.New("didn't get a username")
	}
	uid, err = strconv.Atoi(user.Uid)
	gid, err = strconv.Atoi(user.Gid)
	return user.Username, user.HomeDir, uid, gid, nil
}
