package main

import (
	"os"
	"testing"
)

func Test_loadAllConfig_1(t *testing.T) {
	if _, err := os.Stat("/etc/dockersh"); os.IsNotExist(err) {
		t.Log("No /etc/dockersh, skipping test")
		return
	}
	_, err := loadAllConfig("notexist", "/notexist")
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func Test_Dump_1(t *testing.T) {
	exp := "ImageName busybox MountHomeTo %h ContainerUsername %u Shell /bin/ash DockerSocket /var/run/docker.sock"
	str := defaultConfig.Dump()
	if str != exp {
		t.Errorf("Got '%s' expected '%s'", str, exp)
	}

}

func Test_DefaultConfig_1(t *testing.T) {
	if defaultConfig.ImageName == "busybox" {
		t.Log("default ImageName passed.")
	} else {
		t.Errorf("default ImageName failed: expected busybox got %s", defaultConfig.ImageName)
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
		t.Errorf("Expected ImageName testimage got %s", c.ImageName)
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
		t.Errorf("Expected Shell dockersg got %s", c.Shell)
	}
	if c.ContainerUsername == "bill" {
		t.Log("set ContainerUserName in user config passed.")
	} else {
		t.Errorf("Expected ContainerUserName bill got %s", c.ContainerUsername)
	}
	if c.ImageName == "fredsimage" {
		t.Log("set ImageName in user config passed.")
	} else {
		t.Errorf("Expected ImageName fredsimage got %s", c.ImageName)
	}
}

func Test_IniConfig_2(t *testing.T) {
	c := Configuration{ContainerUsername: "default_contun", ImageName: "default", EnableUserConfig: true, EnableUserImageName: true}
	n, err := loadConfigFromString([]byte(`[dockersh]
imagename = testimage
containerusername = shouldbeblacklisted`), "fred")
	c = mergeConfigs(c, n, true)
	if err != nil {
		t.Error(err)
	}
	if c.ImageName == "testimage" {
		t.Log("set ImageName passed.")
	} else {
		t.Errorf("Expected ImageName testimage got %s", c.ImageName)
	}
	if c.ContainerUsername == "default_contun" {
		t.Log("blacklising worked, value not changed")
	} else {
		t.Error("blacklisting failed")
	}
}

func Test_Config_3(t *testing.T) {
	c := Configuration{ContainerUsername: "default_contun", Shell: "default_shell"}
	c, err := loadConfigFromString([]byte(`[dockersh]
shell = global_default
containerusername = global_default
mounthometo = somewhere
enableusershell
`), "fred")
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
	newc, err := loadConfigFromString([]byte(`[dockersh]
shell = user_value
containerusername = user_value
mounthometo = somewhere_else`), "fred")
	if err != nil {
		t.Error(err)
	}
	c = mergeConfigs(c, newc, true)
	if c.Shell != "user_value" {
		t.Error("Local defaults not applying over global defaults")
	} else {
		t.Log("c.shell not overridden")
	}
	if c.ContainerUsername != "global_default" {
		t.Error("Blacklist of container_username in global config failed")
	}
	if c.MountHomeTo != "somewhere" {
		t.Error("Blacklist mounthome in global config failed")
	}
}

func Test_IniConfig_4(t *testing.T) {
	c, err := loadConfigFromString([]byte(`[dockersh]
containerusername = default_contun
imagename = default
enableuserconfig
enableuserimagename`), "fred")
	newc, err := loadConfigFromString([]byte(`[dockersh]
imagename = testimage
containerusername = shouldbeblacklisted`), "fred")
	if err != nil {
		t.Error(err)
	}
	c = mergeConfigs(c, newc, true)
	if c.ImageName == "testimage" {
		t.Log("set ImageName passed.")
	} else {
		t.Errorf("Expected ImageName testimage got %s", c.ImageName)
	}
	if c.ContainerUsername != "default_contun" {
		t.Error("blacklising disabled, value changed")
	} else {
		t.Log("blacklisting enabled, value has not changed")
	}
}

func Test_IniConfig_5(t *testing.T) {
	c, err := loadConfigFromString([]byte(`[dockersh]
containerusername = default_contun
imagename = default
enableuserconfig
`), "fred")
	newc, err := loadConfigFromString([]byte(`[dockersh]
imagename = testimage
containerusername = shouldbeblacklisted`), "fred")
	if err != nil {
		t.Error(err)
	}
	c = mergeConfigs(c, newc, true)
	if c.ImageName == "default" {
		t.Log("set ImageName passed.")
	} else {
		t.Errorf("Expected ImageName default got %s", c.ImageName)
	}
	if c.ContainerUsername != "default_contun" {
		t.Error("blacklising disabled, value changed")
	} else {
		t.Log("blacklisting enabled, value has not changed")
	}
}

func Test_IniConfig_6(t *testing.T) {
	c, err := loadConfigFromString([]byte(`[dockersh]
imagename = default
enableuserimagename

[user "fred"]
imagename = testimage
`), "fred")
	if err != nil {
		t.Error(err)
	}
	if c.ImageName == "testimage" {
		t.Log("set ImageName in user section when blacklisted in [dockersh] passed.")
	} else {
		t.Errorf("Expected ImageName testimage got %s", c.ImageName)
	}
}
