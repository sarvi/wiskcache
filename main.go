package main

import (
	"argparser"
	"config"
	"exec"
	"fmt"

	//"manifest"
	"cache"
	"utils"
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

	/*
		fmt.Println("Wiskcache Mode -- ", ConfigValues.Mode)
		fmt.Println("Wiskcache Base Dir -- ", ConfigValues.BaseDir)
		fmt.Println("Common Envars from config file -- ", ConfigValues.Envars)
		fmt.Println("Command to be executed -- ", CommandtoExec)
	*/
	/*
		if ConfigValues.ToolIdx != -1 {
			fmt.Println("Tool Match Found")
			fmt.Println("Tool Specific Envars -- ", ConfigValues.Tools[ConfigValues.ToolIdx])
		} else {
			fmt.Println("Tool Match Not Found")
		}
	*/
	if ConfigValues.Mode == "readwrite" || ConfigValues.Mode == "verify" {
		env := utils.GetEnvironMap()
		cmdhash, _ := whash.CommandHash(ConfigValues, env, CommandtoExec)
		manifestFile, _ := cache.FindManifest(ConfigValues, cmdhash)
		if !utils.Exists(manifestFile) && ConfigValues.Mode != "verify" {
			//fmt.Println("Running -- ", CommandtoExec)
			// exitcode, logfile, infiles, outfiles := exec.RunCmd(ConfigValues, cmdhash, CommandtoExec)
			_, _, infiles, outfiles, _ := exec.RunCmd(ConfigValues, cmdhash, CommandtoExec)
			/*
			   fmt.Printf("exit: %v\n", exitcode)
			   fmt.Printf("logfile: %v\n", logfile)
			   fmt.Printf("infiles: %v\n", infiles)
			   fmt.Printf("outfiles: %v\n", outfiles)
			*/
			// fmt.Println(exitcode, logfile, infiles, outfiles)
			cache.Create(ConfigValues, infiles, outfiles, manifestFile)
			fmt.Printf("\nCreated manifest: %v, copied output to cache\n", manifestFile)
		} else if !utils.Exists(manifestFile) && ConfigValues.Mode == "verify" {
			fmt.Printf("%v is not found and can't verify\n", manifestFile)
		} else if utils.Exists(manifestFile) && ConfigValues.Mode == "verify" {
			fmt.Println("Verifying ...")
			if cache.Verify(ConfigValues, manifestFile) {
				fmt.Println("All Matched.")
			}
		} else {
			fmt.Printf("Found manifest: %v, copying out from cache ...\n", manifestFile)
			cache.CopyOut(ConfigValues, manifestFile)
			fmt.Println("Done!")
		}
	}

	// args := os.Args[1:]
	// exec.cmdhash(args)
}
