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
	var ConfigValues config.Config
	var CommandtoExec []string

	ConfigValues, CommandtoExec = argparser.ArgParse()
	ConfigValues.ToolIdx = config.ToolMatcher(ConfigValues, CommandtoExec)
	if ConfigValues.Mode == "learn" {
		ConfigValues.CacheBaseDir = filepath.Join(ConfigValues.BaseDir, "LearningCache")
	}

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
