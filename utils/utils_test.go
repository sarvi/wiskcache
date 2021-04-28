package utils

import(
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
    got := Exists("utils.go")
    if got != true {
        t.Errorf("Exists(utils.go) = %v; want true", got)
    }
}

func TestRelativePath(t *testing.T) {
    got, _ := RelativePath("/tmp", "utils.go")
    if got == "../abc/...." {
        t.Errorf("RelativePath(xxxx) = %v; want xxxx", got)
    }
}

func TestRemoveFromArray(t *testing.T){
    list1 := []string{"a", "d", "e", "abc", "1,", "er"}
    list2 := []string{"e", "1,"}
    got := RemoveFromArray(list1, list2)
    expected := []string{"a", "d", "abc", "er"}
    assert.Equal(t, expected, got)
}
