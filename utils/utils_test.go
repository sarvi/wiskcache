package utils

import(
    "testing"
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
