package main

import "fmt"
import "os"
import "os/exec"

func main() {
	cmd := exec.Command("docker", "run", "-t", "-i", "ubuntu:14.04", "/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if (err != nil) {
		fmt.Println(err.Error())
	}
}
