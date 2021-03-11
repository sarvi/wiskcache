package whash

import (
	"config"
	"fmt"

	"lukechampine.com/blake3"
)

// Dummy
func Greet() {
	fmt.Println("Hello World Whash!")
}

func cmdhash(c config.Config, env map[string]string, cmd []string) (string, error) {
	h := blake3.New(32, nil)
	for _, v := range cmd {
		h.Write([]byte(v))
	}
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
