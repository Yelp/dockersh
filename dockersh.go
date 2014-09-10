package main

import (
	"bufio"
	"fmt"
	"github.com/docker/libcontainer/user"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if os.Args[0] == "/init" {
		os.Exit(initMain())
	} else {
		os.Exit(realMain())
	}
}

func tmplConfigVar(template string, v *configInterpolation) string {
	shell := "/bin/bash"
	r := strings.NewReplacer("%h", v.Home, "%u", v.User, "%s", shell) // Arguments are old, new ...
	return r.Replace(template)
}

func getInterpolatedConfig(config *Configuration, configInterpolations configInterpolation) error {
	config.ContainerUsername = tmplConfigVar(config.ContainerUsername, &configInterpolations)
	config.MountHomeTo = tmplConfigVar(config.MountHomeTo, &configInterpolations)
	config.MountHomeFrom = tmplConfigVar(config.MountHomeFrom, &configInterpolations)
	config.ImageName = tmplConfigVar(config.ImageName, &configInterpolations)
	config.Shell = tmplConfigVar(config.Shell, &configInterpolations)
	config.UserCwd = tmplConfigVar(config.UserCwd, &configInterpolations)
	config.ContainerName = tmplConfigVar(config.ContainerName, &configInterpolations)
	return nil
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func gatewayIP() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", errors.New("Could not open /proc/net/route")
	}
	defer file.Close()
	r := bufio.NewReader(file)
	s, err := Readln(r)
	ip := ""
	for err == nil {
		f := strings.Fields(s)
		if f[1] == "00000000" {
			a, _ := hex.DecodeString(f[2])
			ip = fmt.Srintf("%v.%v.%v.%v", a[3], a[2], a[1], a[0])
			err = nil
			break
		}
		s, err = Readln(r)
	}
	return ip, err
}

func initMain() int {
	fmt.Fprintf(os.Stdout, "started dockersh persistent container\n")
	if file, err := os.Open("/portforward"); err == nil {
		fmt.Printf("/portforward file exists; processing...")
		r := bufio.NewReader(file)
		s, err := Readln(r)
		for err == nil {
			err = validatePortforwardString(s)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
			parts := strings.Split(s, ":") // Parts is hostport:containerport
			localAddr := "127.0.0.1:" + parts[1]
			gw, err := gatewayIP()
			if err != nil {
				panic(err)
			}
			remoteAddr := gw + ":" + parts[0]
			go proxyMain(localAddr, remoteAddr)
			s, err = Readln(r)
		}
		file.Close()
	}
	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc
	return 0
}

func realMain() int {
	err := dockerVersionCheck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Docker version error: %v", err)
		return 1
	}
	username, homedir, uid, gid, err := getCurrentUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v", err)
		return 1
	}
	config, err := loadAllConfig(username, homedir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config: %v\n", err)
		return 1
	}
	configInterpolations := configInterpolation{homedir, username}
	err = getInterpolatedConfig(&config, configInterpolations)
	if err != nil {
		panic(fmt.Sprintf("Cannot interpolate config: %v", err))
	}

	_, err = dockerpid(config.ContainerName)
	if err != nil {
		_, err = dockerstart(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start container: %s\n", err)
			return 1
		}
	}
	_, _, groups, _, err := user.GetUserGroupSupplementaryHome(username, 65536, 65536, "/")
	err = nsenterexec(config.ContainerName, uid, gid, groups, config.UserCwd, config.Shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting shell in new container: %v\n", err)
		return 1
	}
	return 0
}
