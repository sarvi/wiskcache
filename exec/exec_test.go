package exec

import (
	"fmt"
	"testing"
)

func TestIsUpper(t *testing.T) {
	fmt.Println("Test IsUpper Exec")
	Greet()
}

func TestIsLower(t *testing.T) {
	fmt.Println("Test IsLower Exec")
	t.Fail()
}
