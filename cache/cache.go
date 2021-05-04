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
    "io"
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
    manif := manifest.FileManifest{InputFile:[][]string{}, OutputFile:[][]string{}}
    for{
        if !utils.Exists(manifestFile){
            return manifestFile, nil
        }else{
            mismatch = false
            manifestFromFile, _ := manifest.ReadManifest(manifestFile)
            
            for _, inputFile := range manifestFromFile.InputFile{
                fullpath := inputFile[0]
                if !filepath.IsAbs(fullpath){
                    fullpath = filepath.Join(config.BaseDir, fullpath)
                }
                if !utils.Exists(fullpath){
                    manifestFile = filepath.Join(cacheDir, fmt.Sprintf("manifest.%v", inputFile[1]))
                    return manifestFile, nil
                }else{
                    hash, _ := manifest.GetHash(fullpath)
                    manif.InputFile = append(manif.InputFile, []string{inputFile[0], hash})
                }
            }
            for inputIndex, inputFile := range manifestFromFile.InputFile{
                if manif.InputFile[inputIndex][1] != inputFile[1]{
                    manifestFile = filepath.Join(cacheDir, fmt.Sprintf("manifest.%v", inputFile[1]))
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

func Create(config config.Config, logFile string, inFile []string, outFile []string, symLinks [][2]string, manifestfile string)(error){
    // create manifest file and copy outputfiles to cache

    // manifestfile is retrieved from FindManifest
    // create manifest file
    infile, _ := utils.ConverFilesToRelativePath(config, inFile)
    outfile, _ := utils.ConverFilesToRelativePath(config, outFile)
    err := manifest.SaveManifestFile(config, logFile, infile, outfile, symLinks, manifestfile)
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
    done := make(chan bool, 1)
    go func(){
        for _, ofile := range outfile{
            // full path means it's not a file in workspace
            if filepath.IsAbs(ofile){
                continue
            }
            target := filepath.Join(dirOfCachedOutputFiles, strings.Replace(ofile, "/", ".", -1))
            fmt.Printf("Copying %v to %v", filepath.Join(config.BaseDir, ofile), target)
            cpCmd := exec.Command("cp", filepath.Join(config.BaseDir, ofile), target)
            err = cpCmd.Run()
            if err != nil{
                break
            }
        }
        done <- true
    }()
    <-done
    if err != nil{
        return err
    }

    go func(){
        if logFile != "" && utils.Exists(logFile){
            cpCmd := exec.Command("cp", logFile,
                                  filepath.Join(dirOfCachedOutputFiles, filepath.Base(logFile)))
            err = cpCmd.Run()
        }
        done <- true
    }()
    <-done
    if err != nil{
        return err
    }

    return nil
}

func CopyOut(config config.Config, manifestFile string)(error){
    // copy from cache
    var err error
    dirOfCachedOutputFiles := filepath.Join(filepath.Dir(manifestFile),
                                            strings.Replace(filepath.Base(manifestFile), "manifest.", "", 1))
    manifestdata, _ := manifest.ReadManifest(manifestFile)
    done := make(chan bool, 1)
    go func(){
        for _, outputFile := range manifestdata.OutputFile{
            // if outputFile is abs path, it's not a file in workspace then
            if filepath.IsAbs(outputFile[0]){
                continue
            }
            srcFile := filepath.Join(dirOfCachedOutputFiles, strings.Replace(outputFile[0], "/", ".", -1))
            tgtFile := filepath.Join(config.BaseDir, outputFile[0])
            dirOfTgt := filepath.Dir(tgtFile)
            if !utils.Exists(dirOfTgt){
                err = os.MkdirAll(dirOfTgt, 0775)
                if err != nil{
                    break
                }
            }
            fmt.Printf("Copying %v to %v\n", srcFile, tgtFile)
            cpCmd := exec.Command("cp", srcFile, tgtFile)
            err = cpCmd.Run()
            if err != nil{
                break
            }
        }
        done <- true
    }()
    <- done
    if err != nil{
        return err
    }

    go func(){
        for _, symLink := range manifestdata.SymLink{
            os.Symlink(symLink[1], filepath.Join(config.BaseDir, symLink[0]))
        }
        done <- true
    }()
    <- done

    // print out log file
    if manifestdata.LogFile != ""{
        logFile := filepath.Join(dirOfCachedOutputFiles, manifestdata.LogFile)
        if utils.Exists(logFile){
            file, err := os.Open(logFile)
            if err != nil{
                fmt.Printf("Failed to open %v\n", logFile)
            }else{
                _, err = io.Copy(os.Stdout, file)
                if err != nil {
                    fmt.Printf("io.Copy failed  %v\n", err)
                }
            }
        }
    }
    return nil
}

func Verify(config config.Config, manifestFile string)(bool){
    manifestdata, _ := manifest.ReadManifest(manifestFile)
    matched := true
    for _, outputFile := range manifestdata.OutputFile{
        fullpath := outputFile[0]
        hash := outputFile[1]
        if !filepath.IsAbs(fullpath){
            fullpath = filepath.Join(config.BaseDir, fullpath)
        }
        hashOfFileInWorkspace, _ := manifest.GetHash(fullpath)
        fmt.Printf("Comparing %v ...\n", outputFile)
        if hash != hashOfFileInWorkspace{
            fmt.Printf("%v is not matched, hash: %v, hashInWorkspace %v\n", fullpath, hash, hashOfFileInWorkspace)
            matched = false 
        }
    }
    return matched
}
