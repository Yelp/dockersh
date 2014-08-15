package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-namespaces/namespace"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/namespaces"
	"github.com/docker/libcontainer/security/capabilities"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	. "syscall"
)

func loadContainer(path string) (*libcontainer.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var container *libcontainer.Config
	if err := json.NewDecoder(f).Decode(&container); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()
	return container, nil
}

func openNamespaceFd(pid int, path string) (*os.File, error) {
	return os.Open(fmt.Sprintf("/proc/%s/root%s", strconv.Itoa(pid), path))
}

func dropCaps() (err error) {
	// Drop capabilities except those in the whitelist, from https://github.com/docker/docker/blob/master/daemon/execdriver/native/template/default_template.go
	cape := capabilities.DropBoundingSet([]string{
		"CHOWN",
		"DAC_OVERRIDE",
		"FSETID",
		"FOWNER",
		//"MKNOD",
		//"NET_RAW",
		//"SETGID",
		//"SETUID",
		"SETFCAP",
		"SETPCAP",
		"NET_BIND_SERVICE",
		"SYS_CHROOT",
		"KILL",
		"AUDIT_WRITE",
	})
	if cape != nil {
		panic(cape)
	}
	return nil
}

func nsenterexec(containerName string, uid int, gid int, groups []int, wd string, shell string) (err error) {
	containerpid, err := dockerpid(containerName)
	if err != nil {
		panic(fmt.Sprintf("Could not get PID for container: %s", containerName))
	}
	containerSha, err := dockersha(containerName)
	if err != nil {
		panic(fmt.Sprintf("Could not get SHA for container: %s %s", err.Error(), containerName))
	}
	containerConfigLocation := fmt.Sprintf("/var/lib/docker/execdriver/native/%s/container.json", containerSha)
	container, err := loadContainer(containerConfigLocation)
	if err != nil {
		panic(fmt.Sprintf("Could not load container configuration: %v", err))
	}

	rootfd, err := openNamespaceFd(containerpid, "")
	if err != nil {
		panic(fmt.Sprintf("Could not open fd to root: %s", err))
	}
	rootfd.Close()

	cwdfd, err := openNamespaceFd(containerpid, wd)
	if err != nil {
		panic(fmt.Sprintf("Could not open fs to working directory (%s): %s", wd, err))
	}
	cwdfd.Close()

	if strings.HasPrefix(shell, "/") != true {
		return errors.New(fmt.Sprintf("Shell '%s' does not start with /, need an absolute path", shell))
	}
	shell = path.Clean(shell)
	shellfd, err := openNamespaceFd(containerpid, shell)
	shellfd.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot find your shell %s inside your container", shell))
	}

	var nslist = []uintptr{namespace.CLONE_NEWIPC, namespace.CLONE_NEWUTS, namespace.CLONE_NEWNET, namespace.CLONE_NEWPID, namespace.CLONE_NEWNS} // namespace.CLONE_NEWUSER
	for _, ns := range nslist {
		nsfd, err := namespace.OpenProcess(containerpid, ns)
		if nsfd == 0 || err != nil {
			panic("namespace.OpenProcess(containerpid, xxx)")
		}
		namespace.Setns(nsfd, ns)
		namespace.Close(nsfd)
	}
	dropCaps()

	pid, err := ForkExec(shell, []string{"sh"}, &ProcAttr{
		//Env:
		Dir: wd,
		//sys.Setsid
		//sys.Setpgid
		//sys.Setctty && sys.Ctty
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys: &SysProcAttr{
			Chroot:     fmt.Sprintf("/proc/%s/root", strconv.Itoa(containerpid)),
			Credential: &Credential{Uid: uint32(uid), Gid: uint32(gid)}, //, Groups: []uint32(groups)},
		},
	})
	if err != nil {
		panic(err)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		panic(fmt.Sprintf("Could not get proc for pid %s", strconv.Itoa(pid)))
	}
	// FIXME Race condition
	cleaner, err := namespaces.SetupCgroups(container, pid)
	if err != nil {
		proc.Kill()
		proc.Wait()
		panic(fmt.Sprintf("SetupCgroups failed: %s", err.Error()))
	}
	if cleaner != nil {
		defer cleaner.Cleanup()
	}

	// err = writeUserMappings(pid, []IdMap{{ContainerId: 0, HostId: uint32(uid)}}, []IdMap{{ContainerId: 0, HostId: uint32(gid)}})
	// if err != nil

	var wstatus WaitStatus
	_, err1 := Wait4(pid, &wstatus, 0, nil)
	if err != nil {
		panic(err1)
	}

	return nil
}

// Stolen from https://raw.githubusercontent.com/mrunalp/libcontainer/152f2faa63f6db55417e84bb4eb52671de820815/forkexec/forkexec.go
type IdMap struct {
	ContainerId uint32
	HostId      uint32
	Size        uint32
}

// Write UID/GID mappings for a process.
func writeUserMappings(pid int, uidMappings, gidMappings []IdMap) error {
	if len(uidMappings) > 5 || len(gidMappings) > 5 {
		return fmt.Errorf("Only 5 uid/gid mappings are supported by the kernel")
	}

	uidMapStr := make([]string, len(uidMappings))
	for i, um := range uidMappings {
		uidMapStr[i] = fmt.Sprintf("%v %v %v", um.ContainerId, um.HostId, um.Size)
	}

	gidMapStr := make([]string, len(gidMappings))
	for i, gm := range gidMappings {
		gidMapStr[i] = fmt.Sprintf("%v %v %v", gm.ContainerId, gm.HostId, gm.Size)
	}

	uidMap := []byte(strings.Join(uidMapStr, "\n"))
	gidMap := []byte(strings.Join(gidMapStr, "\n"))

	uidMappingsFile := fmt.Sprintf("/proc/%v/uid_map", pid)
	gidMappingsFile := fmt.Sprintf("/proc/%v/gid_map", pid)

	if err := ioutil.WriteFile(uidMappingsFile, uidMap, 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(gidMappingsFile, gidMap, 0644); err != nil {
		return err
	}

	return nil
}
