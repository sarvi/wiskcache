package whash

import (
	"config"
	"fmt"
	"testing"
)

func TestHashAgainShouldMatch(t *testing.T) {
	fmt.Println("Test Whash TestHashAgainShouldMatch")
	tconfig := config.Config{
		Envars:  []string{"CWD"},
		Tools:   []config.Tool{},
		ToolIdx: -1,
	}
	tenv := map[string]string{
		"CWD": "/a/b/c",
	}
	tcmd := []string{"g++", "-o", "file.o", "file.c"}
	h1, e1 := cmdhash(tconfig, tenv, tcmd)
	h2, e2 := cmdhash(tconfig, tenv, tcmd)
	fmt.Println(h1, h2)
	if h1 != h2 || e1 != nil || e2 != nil {
		t.Fail()
	}
	Greet()
}

func TestDifferentCWDShouldNotMatch(t *testing.T) {
	fmt.Println("Test Whash TestDifferentCWDShouldNotMatch")
	tconfig := config.Config{
		Envars:  []string{"CWD"},
		Tools:   []config.Tool{},
		ToolIdx: -1,
	}
	tenv1 := map[string]string{
		"CWD": "/a/b/c",
	}
	tenv2 := map[string]string{
		"CWD": "/a/b/x",
	}
	tcmd := []string{"g++", "-o", "file.o", "file.c"}
	h1, e1 := cmdhash(tconfig, tenv1, tcmd)
	h2, e2 := cmdhash(tconfig, tenv2, tcmd)
	fmt.Println(h1, h2)
	if h1 == h2 || e1 != nil || e2 != nil {
		t.Fail()
	}
	Greet()
}
