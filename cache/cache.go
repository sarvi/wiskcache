package cache

import(
    "fmt"
    "config"
    "path/filepath"
    "manifest"
    "sort"
    "utils"
    "os"
    "os/exec"
    "strings"
)

func Greet() {
    fmt.Println("Hello World Cache!")
}

func FindManifest(config config.Config, cmdhash string, infile []string)(string, error){
    // return manifestFile which could exist or not
    cacheDir := filepath.Join(config.CacheBaseDir, cmdhash)
    if !utils.Exists(cacheDir){
        err := os.MkdirAll(cacheDir, 0775)
        if err != nil{
            return "", err
        }
    }
    manif := manifest.GenerateManifest(infile, []string{}, config.BaseDir)
    manifestFile := filepath.Join(cacheDir, "manifest.base")
    updManifestFile := false
    mismatch := false
    for{
        if !utils.Exists(manifestFile) || updManifestFile == true{
            return manifestFile, nil
        }else{
            mismatch = false
            manifestFromFile, _ := manifest.ReadManifest(manifestFile)
            if manif.InputFile == nil{
                updManifestFile = true
                continue
            }
            // sort InputFile by file name
            keyOfInputFile := make([]string, 0, len(manif.InputFile))
            for key := range manif.InputFile{
                keyOfInputFile = append(keyOfInputFile, key)
            }
            sort.Strings(keyOfInputFile)
            for _, inputfile := range keyOfInputFile{
                if manif.InputFile[inputfile] != manifestFromFile.InputFile[inputfile]{
                    manifestFile = filepath.Join(cacheDir, fmt.Sprintf("manifest.%v", manif.InputFile[inputfile]))
                    mismatch = true
                    break
                }
            }
            if mismatch == false{
                return manifestFile, nil
            }
        }
    }
    return "", nil
}

func Create(config config.Config, infile []string, outfile []string, manifestfile string)(error){
    // create manifest file and copy outputfiles to cache

    // manifestfile is retrieved from FindManifest
    // create manifest file
    err := manifest.SaveManifestFile(config, infile, outfile, manifestfile)
    if err != nil{
        return err
    }

    // copy outputfiles to cache
    dirOfCachedOutputFiles := filepath.Join(filepath.Dir(manifestfile),
                                            strings.Replace(filepath.Base(manifestfile), "manifest.", "", 1))
    if !utils.Exists(dirOfCachedOutputFiles){
        err = os.MkdirAll(dirOfCachedOutputFiles, 0775)
        if err != nil{
            return err
        }
    }
    for _, ofile := range outfile{
        fullPath := filepath.Join(dirOfCachedOutputFiles, filepath.Dir(ofile))
        if !utils.Exists(fullPath){
            err = os.MkdirAll(fullPath, 0775)
            if err != nil{
                return err
            }
        }
        cpCmd := exec.Command("cp", filepath.Join(config.BaseDir, ofile),
                              filepath.Join(fullPath, filepath.Base(ofile)))
        err = cpCmd.Run()
        if err != nil{
            return err
        }
    }
    return nil
}

func CopyOut(config config.Config, manifestFile string)(error){
    // copy from cache
    var err error
    dirOfCachedOutputFiles := filepath.Join(filepath.Dir(manifestFile),
                                            strings.Replace(filepath.Base(manifestFile), "manifest.", "", 1))
    manifestdata, _ := manifest.ReadManifest(manifestFile)
    for outputFile, _ := range manifestdata.OutputFile{
        srcFile := filepath.Join(dirOfCachedOutputFiles, outputFile)
        tgtFile := filepath.Join(config.BaseDir, outputFile)
        dirOfTgt := filepath.Dir(tgtFile)
        if !utils.Exists(dirOfTgt){
            err = os.MkdirAll(dirOfTgt, 0775)
            if err != nil{
                return err
            }
        }
        cpCmd := exec.Command("cp", srcFile, tgtFile)
        err = cpCmd.Run()
        if err != nil{
            return err
        }
    }
    return nil

}
