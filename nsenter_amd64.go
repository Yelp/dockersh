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

const (
	CLONE_VFORK = 0x00004000 /* set if the parent wants the child to wake it up on mm_release */
	SIGCHLD     = 0x14       /* Should set SIGCHLD for fork()-like behavior on Linux */
)

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

	// see go/src/pkg/syscall/exec_unix.go
	syscall.ForkLock.Lock()

	// Stolen from https://github.com/tobert/lnxns/blob/master/src/lnxns/nsfork_linux.go
	r1, _, err1 := syscall.RawSyscall(syscall.SYS_CLONE, uintptr(CLONE_VFORK), 0, 0)

	syscall.ForkLock.Unlock()

	if err1 != 0 {
		panic(err1)
	}

	// parent will get the pid, child will be 0
	if int(r1) != 0 {
		// Parent
		proc, _ := os.FindProcess(int(pid))
		proc.Wait()
		return nil
	}

	// Child

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
