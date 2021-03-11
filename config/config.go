package config

import "fmt"

type Tool struct {
	Match  string
	Envars []string
}

// declaring a student struct
type Config struct {
	// declaring struct variables
	BaseDir string
	Envars  []string
	Tools   []Tool
	ToolIdx int
}

func Greet() {
	fmt.Println("Hello World Config!")
}
