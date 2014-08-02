package main

import "testing"
import "fmt"

func Test_DefaultConfig_1(t *testing.T) {
    if defaultConfig.ImageName == "ubuntu" {
		t.Log("default ImageName passed.")
	} else {
		t.Error(fmt.Sprintf("default ImageName failed: expected ubuntu got %s", defaultConfig.ImageName))
	}
}
