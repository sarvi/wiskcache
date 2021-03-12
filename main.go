package main

import (
	"argparser"
	"config"
	"exec"
	"fmt"
	"manifest"
	"whash"
)

func main() {
	// env := os.Environ()
	// conf, cmd = argparser.ArgParse()
	// cmdhash = whash.cmdhash(conf, env, cmd)
	// man = cache.findmanifest(cmdhash)
	// if config.mode == "active" || config.mode == "verify" {
	// 	if man == nil {
	// 		exitcode, infiles, outfiles = exec.runcmd(cmd)
	// 		if exitcode != 0 {
	// 			exit(exitcode)
	// 		}
	// 		man = manifest.new(man, infile, outfiles)
	// 		cache.create(man)
	// 		exit(0)
	// 	} else {
	// 		if config.mode == "active" {
	// 			cache.copyout(man) // Copy out content from cache
	// 		} else {
	// 			cache.verify(man) // Verify content of Manifest with content of workspace
	// 		}
	// 		exit(0)
	// 	}
	// } else {
	// 	manfiest.learn(conf, env, cmd, cmdhash)
	// }

	var ConfigValues config.Config
	var CommandtoExec []string

	ConfigValues, CommandtoExec = argparser.ArgParse()
	ConfigValues.ToolIdx = config.ToolMatcher(ConfigValues, CommandtoExec)

	fmt.Println("Wiskcache Mode -- ", ConfigValues.Mode)
	fmt.Println("Wiskcache Base Dir -- ", ConfigValues.BaseDir)
	fmt.Println("Common Envars from config file -- ", ConfigValues.Envars)
	fmt.Println("Command to be executed -- ", CommandtoExec)
	if ConfigValues.ToolIdx != -1 {
		fmt.Println("Tool Match Found")
		fmt.Println("Tool Specific Envars -- ", ConfigValues.Tools[ConfigValues.ToolIdx])
	} else {
		fmt.Println("Tool Match Not Found")
	}
	env := map[string]string{}
	cmdhash, _ := whash.CommandHash(ConfigValues, env, CommandtoExec)
	exitcode, logfile, infiles, outfiles := exec.RunCmd(ConfigValues, cmdhash, CommandtoExec)
	fmt.Println(exitcode, logfile, infiles, outfiles)

	hash, err := manifest.GetHash("main.go")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hash)
	}
	// args := os.Args[1:]
	// exec.cmdhash(args)
}
