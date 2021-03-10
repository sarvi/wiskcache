package manifest

import (
    "os"
    "io"
    "fmt"
    "errors"
    "utils"
    "lukechampine.com/blake3"
)

type FileManifest struct{
    inputFile map[string]string
    outputFile map[string]string
}

func GetHash(file string)(string, error){
    if utils.Exists(file){
        f, err := os.Open(file)
        if err != nil {
            return "", err
        }
        defer f.Close()
        h := blake3.New(32, nil)
        if _, err := io.Copy(h, f); err != nil {
            return "", err
        }
        return fmt.Sprintf("%x", h.Sum(nil)), nil
    }else{
        return "", errors.New(fmt.Sprintf("%v does not exist", file))
    }
}

func MatchHash(file string, hash string)(string, error){
    if utils.Exists(file){
        new_hash, err := GetHash(file) 
        if err != nil {
            return "", err
        }
        if hash == new_hash {
            return "", nil
        }else{
            return new_hash, nil
        }
    }else{
        return "", errors.New(fmt.Sprintf("%v does not exist", file))
    }
}

func GenerateManifest(inputFileList []string, outputFileList []string)(FileManifest){
    manifest := FileManifest{inputFile:make(map[string]string), outputFile:make(map[string]string)}
    for _, file := range inputFileList{
        hash, err := GetHash(file)
        if err == nil{
            manifest.inputFile[file] = hash
        }
    } 
    for _, file := range outputFileList{
        hash, err := GetHash(file)
        if err == nil{
            manifest.outputFile[file] = hash
        }
    } 
    return manifest
}
