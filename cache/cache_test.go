package cache

import (
	"fmt"
	"testing"
)

func TestIsUpper(t *testing.T) {
	fmt.Println("Test IsUpper Cache")
	Greet()
}

func TestIsLower(t *testing.T) {
	fmt.Println("Test IsLower Cache")
	t.Fail()
}
