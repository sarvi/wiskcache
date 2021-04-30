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
    InputFile [][]string `json:"inputfile"`
    OutputFile [][]string `json:"outputfile"`
    SymLink [][]string `json:"symlink"`
    LogFile string `json:"logfile"`
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

func GenerateManifest(logFile string, inputFileList []string, outputFileList []string, symLinks [][2]string, baseDirOfWorkspace string)(FileManifest){
    manifest := FileManifest{InputFile:[][]string{}, OutputFile:[][]string{}, SymLink:[][]string{}, LogFile:""}
    manifest.LogFile = filepath.Base(logFile)
    for _, file := range inputFileList{
        fullpath := file
        if !filepath.IsAbs(fullpath){ 
            fullpath = filepath.Join(baseDirOfWorkspace, file)
        }
        hash, err := GetHash(fullpath)
        if err == nil{
            manifest.InputFile = append(manifest.InputFile, []string{file, hash})
        }
    } 
    for _, file := range outputFileList{
        fullpath := file
        if !filepath.IsAbs(fullpath){ 
            fullpath = filepath.Join(baseDirOfWorkspace, file)
        }
        hash, err := GetHash(fullpath)
        if err == nil{
            manifest.OutputFile = append(manifest.OutputFile, []string{file, hash})
        }
    } 
    for _, symlink := range symLinks{
        manifest.SymLink = append(manifest.SymLink, []string{symlink[0], symlink[1]})
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

func SaveManifestFile(config config.Config, logFile string, inputFileList []string, outputFileList []string, symLinks [][2]string, manifestFile string)(error){
    // manifestFile is retrieved from cache.FindManifest

    var err error
    // if an output file is in inputFileList as well, remove it from inputFileList
    inputFileList = utils.RemoveFromArray(inputFileList, outputFileList)

    manifest := GenerateManifest(logFile, inputFileList, outputFileList, symLinks, config.BaseDir)
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
