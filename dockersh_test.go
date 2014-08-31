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
