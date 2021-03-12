package utils

import(
    "os"
    "errors"
    "path/filepath"
    "fmt"
    "config"
    "strings"
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

func ConverFilesToRelativePath(config config.Config, infile []string)([]string, error){
    var err error
    outfile := make([]string, len(infile))
    for i := 0; i < len(infile); i++{
        // a workaround
        infile[i] = strings.Replace(infile[i], "\"", "", -1)
        if filepath.IsAbs(infile[i]) && strings.HasPrefix(infile[i], config.BaseDir) {
            outfile[i], err = RelativePath(config.BaseDir, infile[i])
        }else{
            outfile[i] = infile[i]
        }
    }
    return outfile, err
}
