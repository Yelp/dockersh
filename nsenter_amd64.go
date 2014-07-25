package main

import (
	"github.com/coreos/go-namespaces/namespace"
)
import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func nsenterdetect() (found bool, err error) {
	cmd := exec.Command("/usr/local/bin/nsenter")
	err = cmd.Run()
	if err == nil {
		return true, nil
	}
	/* TODO: Figure out how to get the actual error code from here */
	if e, ok := err.(*exec.ExitError); ok && strings.HasSuffix(e.String(), "1") {
		return false, nil
	}
	return false, err
}

func nsenterexec(pid int, uid int, gid int, wd string, shell string) (err error) {
	// sudo nsenter --target "$PID" --mount --uts --ipc --net --pid --setuid $DESIRED_UID --setgid $DESIRED_GID --wd=$HOMEDIR -- "$REAL_SHELL"
	//cmd := exec.Command("sudo", "/usr/local/bin/nsenter",
	//	"--target", strconv.Itoa(pid), "--mount", "--uts", "--ipc", "--net", "--pid",
	//	"--setuid", strconv.Itoa(uid), "--setgid", strconv.Itoa(gid), fmt.Sprintf("--wd=%s", wd),
	//	"--", shell)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//err = cmd.Run()
	//return err

	/* FIXME: Make these an array and loop through them, as this is gross */

	/* --ipc */
	fd, err := namespace.OpenProcess(pid, namespace.CLONE_NEWIPC)
	if fd == 0 || err != nil {
		panic("namespace.OpenProcess(pid, namespace.CLONE_NEWIPC)")
	}
	namespace.Setns(fd, namespace.CLONE_NEWIPC)
	defer namespace.Close(fd)

	/* --uts */
	fd, err = namespace.OpenProcess(pid, namespace.CLONE_NEWUTS)
	if fd == 0 || err != nil {
		panic("namespace.OpenProcess(pid, namespace.CLONE_NEWUTS)")
	}
	namespace.Setns(fd, namespace.CLONE_NEWUTS)
	defer namespace.Close(fd)

	/* --net */
	fd, err = namespace.OpenProcess(pid, namespace.CLONE_NEWNET)
	if fd == 0 || err != nil {
		panic("namespace.OpenProcess(pid, namespace.CLONE_NEWNET)")
	}
	namespace.Setns(fd, namespace.CLONE_NEWNET)
	defer namespace.Close(fd)

	/* --pid */
	fd, err = namespace.OpenProcess(pid, namespace.CLONE_NEWPID)
	if fd == 0 || err != nil {
		panic("namespace.OpenProcess(pid, namespace.CLONE_NEWPID)")
	}
	namespace.Setns(fd, namespace.CLONE_NEWPID)
	defer namespace.Close(fd)

	/* --mount */
	fd, err = namespace.OpenProcess(pid, namespace.CLONE_NEWNS)
	if fd == 0 || err != nil {
		panic("namespace.OpenProcess(pid, namespace.CLONE_NEWNS)")
	}
	namespace.Setns(fd, namespace.CLONE_NEWNS)
	defer namespace.Close(fd)

	/* END FIXME */

	if gid > 0 {
		err = syscall.Setgroups([]int{}) /* drop supplementary groups */
		if err != nil {
			panic("setgroups failed")
		}
		err = syscall.Setgid(gid)
		if err != nil {
			panic("setgid failed")
		}
	}
	if uid > 0 {
		err = syscall.Setuid(uid)
		if err != nil {
			panic("setuid failed")
		}
	}

	args := []string{shell}
	env := os.Environ()
	execErr := syscall.Exec(shell, args, env)
	if execErr != nil {
		panic(execErr)
	}
	return execErr
}
