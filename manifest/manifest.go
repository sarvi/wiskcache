package manifest

import (
    "os"
    "io"
    "fmt"
    "errors"
    "utils"
    "lukechampine.com/blake3"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "config"
)

type FileManifest struct{
    InputFile map[string]string `json:"inputfile"`
    OutputFile map[string]string `json:"outputfile"`
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

func GenerateManifest(inputFileList []string, outputFileList []string, baseDirOfWorkspace string)(FileManifest){
    manifest := FileManifest{InputFile:make(map[string]string), OutputFile:make(map[string]string)}
    for _, file := range inputFileList{
        fullpath := file
        if !filepath.IsAbs(fullpath){ 
            fullpath = filepath.Join(baseDirOfWorkspace, file)
        }
        hash, err := GetHash(fullpath)
        if err == nil{
            manifest.InputFile[file] = hash
        }
    } 
    for _, file := range outputFileList{
        fullpath := file
        if !filepath.IsAbs(fullpath){ 
            fullpath = filepath.Join(baseDirOfWorkspace, file)
        }
        hash, err := GetHash(fullpath)
        if err == nil{
            manifest.OutputFile[file] = hash
        }
    } 
    return manifest
}

func ReadManifest(manifestFile string)(FileManifest, error){
    var manifest FileManifest
    if !utils.Exists(manifestFile){
        return manifest, errors.New(fmt.Sprintf("%v does not exist.", manifestFile))
    }else{
        data, err := ioutil.ReadFile(manifestFile)
        if err != nil{
            return manifest, errors.New(fmt.Sprintf("Cannot read %v", manifestFile))
        }
        json.Unmarshal(data, &manifest)
        return manifest, nil 
    }
}

func SaveManifestFile(config config.Config, inputFileList []string, outputFileList []string, manifestFile string)(error){
    // manifestFile is retrieved from cache.FindManifest
    var err error
    manifest := GenerateManifest(inputFileList, outputFileList, config.BaseDir)
    jsondata, _ := json.MarshalIndent(manifest, "", " ")
    cacheDir := filepath.Dir(manifestFile)
    if !utils.Exists(cacheDir){
        err = os.MkdirAll(cacheDir, 0775)
        if err != nil{
            return err
        }
    } 
    err = ioutil.WriteFile(manifestFile, jsondata, 0664)
    if err != nil{
        return err
    }
    hashOfManifestfile, _ := GetHash(manifestFile)
    // make symlink
    manifestWithHash := filepath.Join(cacheDir, "manifest." + hashOfManifestfile)
    if !utils.Exists(manifestWithHash){
        os.Rename(manifestFile, manifestWithHash)
    }else{
        os.Remove(manifestFile)
    }
    relativePath, _ := utils.RelativePath(cacheDir, manifestWithHash)
    os.Symlink(relativePath, manifestFile)
    return nil
}
