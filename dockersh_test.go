package main

import (
	"testing"
)

func Test_templConfigVar_1(t *testing.T) {
	i := configInterpolation{"foo", "bar"}
	out := tmplConfigVar("%s", &i)
	if out == "/bin/bash" {
		t.Log("OK")
	} else {
		t.Error("Expected /bin/bash, got %s", out)
	}
}

func Test_getInterpolatedConfig_1(t *testing.T) {
	i := configInterpolation{"foo", "bar"}
	c := defaultConfig
	e := getInterpolatedConfig(&c, i)
	if e != nil {
		t.Error("Error")
	}
	if c.MountHomeFrom != "foo" {
		t.Errorf("MountHomeFrom is %s not foo", c.MountHomeFrom)
	}
}

func Test_gatewayIP_1(t *testing.T) {
	ip, err := gatewayIP()
        if err != nil {
		t.Error("Error")
        }
	t.Log(ip)
}
