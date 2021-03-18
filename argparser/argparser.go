package argparser

import (
	"config"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//ArgParse funtion to parse arguments
func ArgParse() (ConfigValues config.Config, CommandLine []string) {

	var defaultConfigFile string = os.Getenv("WISKCACHE_CONFIG")
	var defaultBaseDir string = os.Getenv("WISK_WSROOT")
	if defaultBaseDir == "" {
		defaultBaseDir, _ = os.Getwd()
	}
	mode := flag.String("mode", "", "Wiskcache operating mode. Possible values - readonly, writeonly, readwrite, learn, verify")
	var InputconfigFile string
	baseDir := flag.String("base_dir", defaultBaseDir, "Wiskcache will rewrite absolute paths beginning with base_dir into paths relative to the current working directory")
	if defaultConfigFile == "" {
		defaultConfigFile = filepath.Join(*baseDir, "wisk/wiskcache_config.yaml")
	}
	flag.StringVar(&InputconfigFile, "config", defaultConfigFile, "Wiskcache configure file location")
	flag.Parse()
	remainingArgs := flag.Args()
	CommandLine = remainingArgs

	//Parsing config file to get Configvalues instance
	if InputconfigFile != "" {
		_, err := os.Stat(InputconfigFile)
		if err == nil {
			ConfigValues = config.Parseconfig(InputconfigFile)
		} else {
			fmt.Println("Config file does not exist")
			os.Exit(1)
		}
	}

	//ToolIndex default value set to -1
	ConfigValues.ToolIdx = -1

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	ConfigValues.UserName = user.Username

	//Adding the Wiskcache Mode and Base directory information to Config instance
	if *mode != "" {
		ConfigValues.Mode = *mode
	} else if ConfigValues.Mode == "" {
		ConfigValues.Mode = "readwrite"
	}

	if ConfigValues.WiskTrackLib == "" {
		libpath, err := os.Executable()
		if err != nil {
			fmt.Println("Cannot Locate Wisk Track Library")
			log.Fatal(err)
		}
		if filepath.HasPrefix(libpath, "/router/bin/") {
			libpath = strings.Replace(libpath, "/router", "/sw/packages/wisk/current", 1)
		}
		libpath = filepath.Join(filepath.Dir(filepath.Dir(libpath)), "${LIB}", "libwisktrack.so")
		ConfigValues.WiskTrackLib = libpath
	}

	if *baseDir != "" {
		ConfigValues.BaseDir = *baseDir
	}

	//Validating arguments
	ModeValues := []string{"readonly", "writeonly", "readwrite", "learn", "verify"}
	if !contains(ModeValues, ConfigValues.Mode) {
		fmt.Println("Invalid mode. Possible values - readonly, writeonly, readwrite, learn, verify")
		os.Exit(1)
	}

	if ConfigValues.CacheBaseDir == "" {
		log.Fatalf("CacheBaseDir is not set in %s", InputconfigFile)
	}

	if len(ConfigValues.Envars) == 0 {
		ConfigValues.Envars = []string{"CWD"}
	}

	return

}
