package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

//Tool specific config values
type Tool struct {
	Match      string   `yaml:"Match"`
	ToolEnvars []string `yaml:"ToolEnvars"`
}

//Config structure declaration
type Config struct {
	ToolIdx int      `yaml:"ToolIdx"`
	Mode    string   `yaml:"Mode"`
	BaseDir string   `yaml:"BaseDir"`
	Envars  []string `yaml:"Envars"`
	Tools   []Tool   `yaml:"Tool"`
}

//Greet function for testing
func Greet() {
	fmt.Println("Hello World Config!")
}

//Parseconfig function to parse the config file
func Parseconfig(InputconfigFile string) Config {
	ConfigFile, err := ioutil.ReadFile(InputconfigFile)
	if err != nil {
		fmt.Printf("Error reading Wiskcache configure file: %s\n", err)
	}
	var ConfigValues Config
	err = yaml.Unmarshal(ConfigFile, &ConfigValues)
	if err != nil {
		fmt.Printf("Error parsing Wiskcache configure file: %s\n", err)
	}
	return ConfigValues
}

//ToolMatcher fuction to match the Tool name with the Tool information in config file
func ToolMatcher(ConfigValues Config, CommandLine string) int {

	var matched bool
	var err error
	var idx int = -1
	var Toolsno int = len(ConfigValues.Tools)
	var Toolname string = strings.Split(CommandLine, " ")[0]
	for i := 0; i < Toolsno; i++ {
		matched, err = regexp.MatchString(Toolname, ConfigValues.Tools[i].Match)
		if matched && err == nil {
			idx = i
			break
		}
	}
	return idx
}
