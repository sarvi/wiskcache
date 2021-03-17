package main

import (
	"argparser"
	"config"
	"exec"
	"fmt"
	"path/filepath"
	"strings"

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
	if ConfigValues.Mode == "learn" {
		ConfigValues.CacheBaseDir = filepath.Join(ConfigValues.BaseDir, "LearningCache")
	}

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
	cmdexeced := false
	infiles := []string{}
	outfiles := []string{}
	manifestFile := ""
	env := utils.GetEnvironMap()
	cmdhash, _ := whash.CommandHash(ConfigValues, env, CommandtoExec)
	if strings.HasPrefix(ConfigValues.Mode, "read") || ConfigValues.Mode == "verify" {
		manifestFile, _ = cache.FindManifest(ConfigValues, cmdhash)
		if ConfigValues.Mode != "verify" {
			fmt.Printf("Found manifest: %v, copying out from cache ...\n", manifestFile)
			cache.CopyOut(ConfigValues, manifestFile)
			fmt.Println("Done!")
		}
	}
	if !utils.Exists(manifestFile) || ConfigValues.Mode == "verify" || ConfigValues.Mode == "learn" {
		cmdexeced = true
		_, _, infiles, outfiles, _ = exec.RunCmd(ConfigValues, cmdhash, CommandtoExec)
	}
	if strings.Contains(ConfigValues.Mode, "write") && cmdexeced {
		cache.Create(ConfigValues, infiles, outfiles, manifestFile)
		fmt.Printf("\nCreated manifest: %v, copied output to cache\n", manifestFile)
	}
	if ConfigValues.Mode == "verify" && cmdexeced {
		if utils.Exists(manifestFile) {
			fmt.Println("Verifying ...")
			if cache.Verify(ConfigValues, manifestFile) {
				fmt.Println("All Matched.")
			}
		} else {
			fmt.Printf("%v is not found and can't verify\n", manifestFile)
		}
	}
	if ConfigValues.Mode == "learn" && cmdexeced {
		fmt.Println("Learn Mode, executing the command second time, collect learning data ...")
		_, _, infiles, outfiles, _ = exec.RunCmd(ConfigValues, cmdhash, CommandtoExec)
	}
}
