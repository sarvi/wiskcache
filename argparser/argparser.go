package argparser

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

//Greet function for testcase
func Greet() {
	fmt.Println("Hello World Argparse!")
}

//Tool structure for tool specific values in config file
type Tool struct {
	ToolEnvars []string `yaml:"ToolEnvars"`
	Match      string   `yaml:"Match"`
}

//CacheConfig -- Cache configure file values
type CacheConfig struct {
	ToolIdx int      `yaml:"ToolIdx"`
	Envars  []string `yaml:"Envars"`
	Tools   []Tool   `yaml:"Tool"`
}

//RunCmd -- To exectue shell commands
func RunCmd(command string) (string, string) {
	cmd := exec.Command("bash", "-c", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	return stdout.String(), stderr.String()
}

//runArgParser funtion to parse arguments
func runArgParser() {

	//parsing the CLI arguments
	mode := flag.String("mode", "active",
		"Wiskcache operating mode. Possible values - active, learning, verify")
	var InputconfigFile string
	flag.StringVar(&InputconfigFile, "config", "<default_config_file_location>", "Wiskcache configure file location")
	baseDir := flag.String("base_dir", "",
		"Wiskcache will rewrite absolute paths beginning with base_dir into paths relative to the current working directory")
	flag.Parse()
	remainingArgs := flag.Args()
	cmdLine := strings.Join(remainingArgs, " ")

	//printing out parsed CLI arguments
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Wiskcache Base Directory -- ", *baseDir)
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Printf("\nWiskcache is operating in '%s' mode\n", *mode)
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Wiskcache configuration file location -- ", InputconfigFile)

	//parsing the input configure file
	ConfigFile, err := ioutil.ReadFile(InputconfigFile)
	if err != nil {
		fmt.Printf("Error reading Wiskcache configure file: %s\n", err)
		return
	}
	var ConfigValues CacheConfig
	err = yaml.Unmarshal(ConfigFile, &ConfigValues)
	if err != nil {
		fmt.Printf("Error parsing Wiskcache configure file: %s\n", err)
	}

	fmt.Println("Configuration file contents -- ")
	//printing out Key-Value pairs of parsed configure file
	fmt.Println("Tool Index -- ", ConfigValues.ToolIdx)
	fmt.Println("Environment Variables -- ", ConfigValues.Envars)

	var Toolsno int
	Toolsno = len(ConfigValues.Tools)

	for i := 0; i < Toolsno; i++ {
		fmt.Println("	Match -- ", ConfigValues.Tools[i].Match)
		fmt.Println("	Tool Environment Variables -- ", ConfigValues.Tools[i].ToolEnvars)
	}
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Command line after Wiskcache flags -- ", cmdLine)
	//fmt.Println("cmd_line_list -- ", remainingArgs)

	//running shell command obtained from CLI and printing output and error
	fmt.Println("Error and output of running the above command line")
	output, errorOut := RunCmd(cmdLine)
	fmt.Println("---stderr---")
	fmt.Println(errorOut)
	fmt.Println("---stdout---")
	fmt.Println(output)
	fmt.Println("--------------------------------------------------------------------------------------------------------")
}
