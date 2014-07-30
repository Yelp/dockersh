package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func getCurrentUser() (username string, homedir string, uid int, gid int, err error) {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v", err)
		return "", "", 0, 0, errors.New("could not get current user")
	}
	return getUser(user)
}

func getUser(user *user.User) (username string, homedir string, uid int, gid int, err error) {
	if user.HomeDir == "" {
		fmt.Fprintf(os.Stderr, "didn't get a home directory")
		return "", "", 0, 0, errors.New("didn't get a home directory")
	}
	if user.Username == "" {
		fmt.Fprintf(os.Stderr, "didn't get a username")
		return "", "", 0, 0, errors.New("didn't get a username")
	}
	uid, err = strconv.Atoi(user.Uid)
	gid, err = strconv.Atoi(user.Gid)
	return user.Username, user.HomeDir, uid, gid, nil
}
