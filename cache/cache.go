package cache

import(
    "fmt"
    "config"
    "path/filepath"
    "manifest"
    "utils"
    "os"
    "os/exec"
    "strings"
)

func Greet() {
    fmt.Println("Hello World Cache!")
}

func FindManifest(config config.Config, cmdhash string)(string, error){
    // return manifestFile which could exist or not
    
    cacheDir := filepath.Join(config.CacheBaseDir, cmdhash)
    if !utils.Exists(cacheDir){
        err := os.MkdirAll(cacheDir, 0775)
        if err != nil{
            return "", err
        }
    }
    manifestFile := filepath.Join(cacheDir, "manifest.base")
    mismatch := false
    manif := manifest.FileManifest{InputFile:[]string{}, InputFileHash:[]string{}, OutputFile:[]string{}, OutputFileHash:[]string{}}
    for{
        if !utils.Exists(manifestFile){
            return manifestFile, nil
        }else{
            mismatch = false
            manifestFromFile, _ := manifest.ReadManifest(manifestFile)
            
            for inputIndex, inputFile := range manifestFromFile.InputFile{ 
                fullpath := inputFile
                if !filepath.IsAbs(fullpath){
                    fullpath = filepath.Join(config.BaseDir, inputFile)
                }
                if !utils.Exists(fullpath){
                    manifestFile = filepath.Join(cacheDir, fmt.Sprintf("manifest.%v", manifestFromFile.InputFileHash[inputIndex]))
                    return manifestFile, nil
                }else{
                    hash, _ := manifest.GetHash(fullpath)
                    manif.InputFile = append(manif.InputFile, inputFile)
                    manif.InputFileHash = append(manif.InputFileHash, hash)
                }
            }
            for inputIndex, _ := range manifestFromFile.InputFile{
                if manif.InputFileHash[inputIndex] != manifestFromFile.InputFileHash[inputIndex]{
                    manifestFile = filepath.Join(cacheDir, fmt.Sprintf("manifest.%v", manifestFromFile.InputFileHash[inputIndex]))
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

func Create(config config.Config, inFile []string, outFile []string, manifestfile string)(error){
    // create manifest file and copy outputfiles to cache

    // manifestfile is retrieved from FindManifest
    // create manifest file
    infile, _ := utils.ConverFilesToRelativePath(config, inFile)
    outfile, _ := utils.ConverFilesToRelativePath(config, outFile)
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
        // full path means it's not a file in workspace
        if filepath.IsAbs(ofile){
            continue
        }
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
    for _, outputFile := range manifestdata.OutputFile{
        // if outputFile is abs path, it's not a file in workspace then
        if filepath.IsAbs(outputFile){
            continue
        }
        srcFile := filepath.Join(dirOfCachedOutputFiles, outputFile)
        tgtFile := filepath.Join(config.BaseDir, outputFile)
        dirOfTgt := filepath.Dir(tgtFile)
        if !utils.Exists(dirOfTgt){
            err = os.MkdirAll(dirOfTgt, 0775)
            if err != nil{
                return err
            }
        }
        fmt.Printf("Copying %v to %v\n", srcFile, tgtFile)
        cpCmd := exec.Command("cp", srcFile, tgtFile)
        err = cpCmd.Run()
        if err != nil{
            return err
        }
    }
    return nil

}

func Verify(config config.Config, manifestFile string)(bool){
    manifestdata, _ := manifest.ReadManifest(manifestFile)
    matched := true
    for outIndex, outputFile := range manifestdata.OutputFile{
        hash := manifestdata.OutputFileHash[outIndex]
        fullpath := outputFile
        if !filepath.IsAbs(outputFile){
            fullpath = filepath.Join(config.BaseDir, outputFile)
        }
        hashOfFileInWorkspace, _ := manifest.GetHash(fullpath)
        fmt.Printf("Comparing %v ...\n", outputFile)
        if hash != hashOfFileInWorkspace{
            fmt.Printf("%v is not matched, hash: %v, hashInWorkspace %v\n", outputFile, hash, hashOfFileInWorkspace)
            matched = false 
        }
    }
    return matched
}
