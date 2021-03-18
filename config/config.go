package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

//Tool specific config values
type Tool struct {
	Match  string   `yaml:"Match"`
	Envars []string `yaml:"Envars"`
}

//Config structure declaration
type Config struct {
	ToolIdx      int      `yaml:"ToolIdx"`
	UserName     string   `yaml:"UserName"`
	Mode         string   `yaml:"Mode"`
	BaseDir      string   `yaml:"BaseDir"`
	WiskTrackLib string   `yaml:"WiskTrackLib"`
	Envars       []string `yaml:"Envars"`
	Tools        []Tool   `yaml:"Tool"`
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
func ToolMatcher(ConfigValues Config, CommandLine []string) (idx int) {

	var matched bool
	var err error
	idx = -1
	var Toolsno int = len(ConfigValues.Tools)
	if len(CommandLine) == 0 {
		log.Fatal("No Command line")
	}
	var Toolname string = CommandLine[0]
	for i := 0; i < Toolsno; i++ {
		matched, err = regexp.MatchString(ConfigValues.Tools[i].Match, Toolname)
		if matched && err == nil {
			idx = i
			break
		}
	}
	return
}
