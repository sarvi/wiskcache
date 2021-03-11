package main

import (
	"cache"
	"config"
	"exec"
	"fmt"
	"manifest"
	"utils"
	"whash"
)

func main() {

	if utils.Exists("go.mod") {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
	cache.Greet()
	exec.Greet()
	whash.Greet()
	config.Greet()
	hash, err := manifest.GetHash("main.go")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hash)
	}
	// args := os.Args[1:]
	// exec.cmdhash(args)
}
