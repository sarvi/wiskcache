package utils

import(
    "os"
    "errors"
    "path/filepath"
    "fmt"
)

func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func RelativePath(basePath string, tgtPath string) (string, error) {
    var relpath = tgtPath
    var err error
    if filepath.IsAbs(tgtPath) {
        relpath, err = filepath.Rel(basePath, tgtPath)
        if err != nil{
            return relpath, err
        }
    }
    // check if path exists
    fullpath := filepath.Join(basePath, relpath)
    if !Exists(fullpath){
        return "", errors.New(fmt.Sprintf("%v does not exist", fullpath))
    }
    return relpath, nil
}