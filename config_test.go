package main

import "testing"
import "fmt"

func Test_DefaultConfig_1(t *testing.T) {
	if defaultConfig.ImageName == "busybox" {
		t.Log("default ImageName passed.")
	} else {
		t.Error(fmt.Sprintf("default ImageName failed: expected busybox got %s", defaultConfig.ImageName))
	}
}

func Test_SimpleConfig_1(t *testing.T) {
	c, err := loadConfigFromString([]byte(``), "fred")
	if err != nil {
		t.Error(err)
	}
	c, err = loadConfigFromString([]byte(`[dockersh]
imagename = testimage`), "fred")
	if err != nil {
		t.Error(err)
	}
	if c.ImageName == "testimage" {
		t.Log("set ImageName passed.")
	} else {
		t.Error(fmt.Sprintf("Expected ImageName testimage got %s", c.ImageName))
	}
}

func Test_UserConfig_1(t *testing.T) {
	c, err := loadConfigFromString([]byte(`[dockersh]
imagename = testimage
shell = someshell

[user "fred"]
imagename = fredsimage
containerusername = bill`), "fred")
	if err != nil {
		t.Error(err)
	}
	if c.Shell == "someshell" {
		t.Log("set Shell in dockersh config passed.")
	} else {
		t.Error(fmt.Sprintf("Expected Shell dockersg got %s", c.Shell))
	}
	if c.ContainerUsername == "bill" {
		t.Log("set ContainerUserName in user config passed.")
	} else {
		t.Error(fmt.Sprintf("Expected ContainerUserName bill got %s", c.ContainerUsername))
	}
	if c.ImageName == "fredsimage" {
		t.Log("set ImageName in user config passed.")
	} else {
		t.Error(fmt.Sprintf("Expected ImageName fredsimage got %s", c.ImageName))
	}
}

/*
func Test_JsonConfig_2(t *testing.T) {
	c := Configuration{BlacklistUserConfig: []string{"container_username"}, ContainerUsername: "default_contun", ImageName: "default"}
	err := loadConfigFromString([]byte(`{"image_name":"testimage","container_username":"shouldbeblacklisted"}`), &c, true)
	if err != nil {
		t.Error(err)
	}
	if c.ImageName == "testimage" {
		t.Log("set ImageName passed.")
	} else {
		t.Error(fmt.Sprintf("Expected ImageName testimage got %s", c.ImageName))
	}
	if c.ContainerUsername == "default_contun" {
		t.Log("blacklising worked, value not changed")
	} else {
		t.Error("blacklisting failed")
	}
}

func Test_JsonConfig_3(t *testing.T) {
	c := Configuration{BlacklistUserConfig: []string{}, ContainerUsername: "default_contun", Shell: "default_shell"}
	err := loadConfigFromString([]byte(`{"shell":"global_default","container_username":"global_default","mount_home_to":"somewhere","blacklist_user_config":"container_username,mount_home_to"}`), &c, true)
	if err != nil {
		t.Error(err)
	}
	if c.Shell != "global_default" {
		t.Error("Set shell to global_default failed")
	}
	if c.ContainerUsername != "global_default" {
		t.Error("Set un to global default failed")
	}
	if c.MountHomeTo != "somewhere" {
		t.Error("Set mounthome to global default failed")
	}
	err = loadConfigFromString([]byte(`{"shell":"user_value","container_username":"user_value","mount_home":"somewhere_else"}`), &c, true)
	if err != nil {
		t.Error(err)
	}
	if c.Shell != "user_value" {
		t.Error("Local defaults not applying over global defaults")
	}
	if c.ContainerUsername != "global_default" {
		t.Error("Blacklist of container_username in global config failed")
	}
	if c.MountHomeTo != "somewhere" {
		t.Error("Blacklist mounthome in global config failed")
	}
}

func Test_JsonConfig_4(t *testing.T) {
	c := Configuration{BlacklistUserConfig: []string{"container_username"}, ContainerUsername: "default_contun", ImageName: "default"}
	err := loadConfigFromString([]byte(`{"image_name":"testimage","container_username":"shouldbeblacklisted"}`), &c, false)
	if err != nil {
		t.Error(err)
	}
	if c.ImageName == "testimage" {
		t.Log("set ImageName passed.")
	} else {
		t.Error(fmt.Sprintf("Expected ImageName testimage got %s", c.ImageName))
	}
	if c.ContainerUsername != "default_contun" {
		t.Log("blacklising disabled, value changed")
	} else {
		t.Error("blacklisting enabled, value has not changes")
	}
}

*/
