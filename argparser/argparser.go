package argparser

import (
	"config"
	"flag"
	"fmt"
	"os"
	"strings"
)

//Greet function for testcase
func Greet() {
	fmt.Println("Hello World Argparse!")
}

//ArgParse funtion to parse arguments
func ArgParse() (config.Config, string) {

	var defaultConfigFile string = os.Getenv("wiskcache_config")
	mode := flag.String("mode", "active", "Wiskcache operating mode. Possible values - active, learning, verify")
	var InputconfigFile string
	flag.StringVar(&InputconfigFile, "config", defaultConfigFile, "Wiskcache configure file location")
	baseDir := flag.String("base_dir", "", "Wiskcache will rewrite absolute paths beginning with base_dir into paths relative to the current working directory")
	flag.Parse()
	remainingArgs := flag.Args()
	CommandLine := strings.Join(remainingArgs, " ")

	var ConfigValues config.Config

	//Parsing config file to get Configvalues instance
	if InputconfigFile != "" {
		ConfigValues = config.Parseconfig(InputconfigFile)
	}

	//ToolIndex default value set to -1
	ConfigValues.ToolIdx = -1

	//Adding the Wiskcache Mode and Base directory information to Config instance
	ConfigValues.Mode = *mode
	ConfigValues.BaseDir = *baseDir

	return ConfigValues, CommandLine

}
