package argparser

import (
	"config"
	"flag"
	"fmt"
	"os"
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
func ArgParse() (ConfigValues config.Config, CommandLine string) {

	var defaultConfigFile string = os.Getenv("WISKCACHE_CONFIG")
	mode := flag.String("mode", "readwrite", "Wiskcache operating mode. Possible values - readonly, writeonly, readwrite, learn, verify")
	var InputconfigFile string
	flag.StringVar(&InputconfigFile, "config", defaultConfigFile, "Wiskcache configure file location")
	baseDir := flag.String("base_dir", "", "Wiskcache will rewrite absolute paths beginning with base_dir into paths relative to the current working directory")
	flag.Parse()
	remainingArgs := flag.Args()
	CommandLine = strings.Join(remainingArgs, " ")

	//Validating arguments
	ModeValues := []string{"readonly", "writeonly", "readwrite", "learn", "verify"}
	if !contains(ModeValues, *mode) {
		fmt.Println("Invalid mode. Possible values - readonly, writeonly, readwrite, learn, verify")
		os.Exit(1)
	}

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

	//Adding the Wiskcache Mode and Base directory information to Config instance
	if *mode != "" {
		ConfigValues.Mode = *mode
	}
	if *baseDir != "" {
		ConfigValues.BaseDir = *baseDir
	}

	return

}
