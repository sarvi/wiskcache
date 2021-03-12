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
	CacheBaseDir string   `yaml:"CacheBaseDir"`
}

//Parseconfig function to parse the config file
func Parseconfig(InputconfigFile string) (ConfigValues Config) {
	ConfigFile, err := ioutil.ReadFile(InputconfigFile)
	if err != nil {
		fmt.Printf("Error reading Wiskcache configure file: %s\n", err)
	}
	err = yaml.Unmarshal(ConfigFile, &ConfigValues)
	if err != nil {
		fmt.Printf("Error parsing Wiskcache configure file: %s\n", err)
	}
	return
}

//ToolMatcher fuction to match the Tool name with the Tool information in config file
func ToolMatcher(ConfigValues Config, CommandLine string) (idx int) {

	var matched bool
	var err error
	idx = -1
	var Toolsno int = len(ConfigValues.Tools)
	var Toolname string = strings.Split(CommandLine, " ")[0]
	for i := 0; i < Toolsno; i++ {
		matched, err = regexp.MatchString(ConfigValues.Tools[i].Match, Toolname)
		if matched && err == nil {
			idx = i
			break
		}
	}
	return
}
