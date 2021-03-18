package whash

import (
	"config"
	"fmt"
	"log"
	"os"
	"strings"

	"lukechampine.com/blake3"
)

func pathtorel(path string, relativeto string) (rv string) {
	rv = strings.Replace(path, relativeto, "", 1)
	if rv != path {
		for strings.Index(rv, "/") == 0 {
			rv = strings.Replace(rv, "/", "", 1)
		}
		rv = strings.Join([]string{"./", rv}, "")
	}
	// fmt.Println("pathtorel: ", path, rv)
	return
}

func cmdnormalize(c config.Config, cmd []string) (rv []string) {
	// fmt.Println("BaseDir: ", c.BaseDir)
	for _, p := range cmd {
		// fmt.Println("replacing: ", p, c.BaseDir)
		rv = append(rv, pathtorel(p, c.BaseDir))
	}
	// fmt.Println("Cmd: ", rv)
	return
}

func envnormalize(c config.Config, env map[string]string) (rv map[string]string) {
	rv = make(map[string]string)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error geting Current Directory")
	} else {
		rv["CWD"] = pathtorel(cwd, c.BaseDir)
	}
	for k, v := range env {
		// fmt.Println("replacing: ", p, c.BaseDir)
		rv[k] = pathtorel(v, c.BaseDir)
	}
	// fmt.Println("Env: ", rv["CWD"])
	return
}

func CommandHash(c config.Config, env map[string]string, cmd []string) (string, error) {
	h := blake3.New(32, nil)
	cmd = cmdnormalize(c, cmd)
	for _, v := range cmd {
		h.Write([]byte(v))
	}
	env = envnormalize(c, env)
	tohashvars := map[string]string{}
	for _, v := range c.Envars {
		if v, exists := env[v]; exists {
			tohashvars[v] = v
		}
	}
	if c.ToolIdx >= 0 {
		for _, v := range c.Tools[c.ToolIdx].Envars {
			if v, exists := env[v]; exists {
				tohashvars[v] = v
			}
		}
	}
	for k, v := range tohashvars {
		h.Write([]byte(k))
		h.Write([]byte(v))
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
