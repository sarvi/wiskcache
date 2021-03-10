package exec

import (
	"fmt"
	"testing"
)

func TestIsUpper(t *testing.T) {
	fmt.Println("Test IsUpper")
	Greet()
}

func TestIsLower(t *testing.T) {
	fmt.Println("Test IsLower")
	t.Fail()
}
