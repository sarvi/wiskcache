package cache

import (
	"fmt"
	"testing"
	"sync"
	"os/exec"
	"strings"
	"os"
	"config"
	"path/filepath"
	"utils"
	"io/ioutil"
)

func TestIsUpper(t *testing.T) {
	fmt.Println("Test IsUpper Cache")
	Greet()
}

func TestIsLower(t *testing.T) {
	fmt.Println("Test IsLower Cache")
	//t.Fail()
}

func TestCopyOutInParallel(t *testing.T) {
	cmd := exec.Command("fallocate", "-l", "1G", "/tmp/1gfile1")
	err := cmd.Run()
        if err != nil{
		fmt.Println("Failed to create files for testing")
		return
        }
	cmd = exec.Command("fallocate", "-l", "1G", "/tmp/1gfile2")
	err = cmd.Run()
        if err != nil{
		fmt.Println("Failed to create files for testing")
		return
        }
	var wg sync.WaitGroup
	outfile := []string{"/tmp/1gfile1", "/tmp/1gfile2"}
	for _, ofile := range outfile{
		target := ofile + ".dest"
		fmt.Printf("Copying %v to %v\n", ofile, target)
		wg.Add(1)
		go func(src string, tgt string){
			defer wg.Done()
			cpCmd := exec.Command("cp", src, tgt)
			cperr := cpCmd.Run()
			if cperr != nil{
	            		err = cperr
			}
		}(ofile, target)
	}
	//cmd = exec.Command("ps", "-ef", "\\|", "grep", "cp")
	cmd = exec.Command("ps", "-ef")
	output, _ := cmd.CombinedOutput()
	wg.Wait()
	os.Remove("/tmp/1gfile1")
	os.Remove("/tmp/1gfile2")
	os.Remove("/tmp/1gfile1.dest")
	os.Remove("/tmp/1gfile2.dest")
	if !strings.Contains(string(output), "cp /tmp/1gfile1 /tmp/1gfile1.dest") || !strings.Contains(string(output), "cp /tmp/1gfile2 /tmp/1gfile2.dest"){
		t.Fail()
	} 
}

func TestManifest(t *testing.T) {
	var config config.Config
	config.CacheBaseDir = "/tmp/cache_testTestFindManifest"
	cmdhash := "cmdhash001"
	manifestFile, _ := FindManifest(config, cmdhash)
	infiles := []string{"../test/hello.c", "../test/square.h", "../test/sum.h"}
        outfiles := []string{"../test/hello.o"}
        symlinks := [][2]string{}
	if filepath.Base(manifestFile) == "manifest.base" {
		cmd := exec.Command("touch", "../test/hello.o")
		cmd.Run()
		// test create a new manifest file
		manifestFile, _ = Create(config, "", infiles, outfiles, symlinks, manifestFile)
		if manifestFile != "/tmp/cache_testTestFindManifest/cmdhash001/manifest.403e4ce21cc1a8d2b887f3255db4559a9aa350034ac035ec90dd46e36c42e792" {
			os.RemoveAll("/tmp/cache_testTestFindManifest")
			t.Fail()
                }
		// test if partial.xxx/manifest.<hashofsublistfilesandhash> symlinks to parent dir's manifest file
		if utils.Exists("/tmp/cache_testTestFindManifest/cmdhash001/partial.1feb0822dbb0ef198ae7453b23513045bd677ef4b4d80e48b10983c9618fa598/manifest.111ac4231fc677d05a106ea0dac0ea5b628b1e8a8e0af287e266cf11f7efe94c"){
			manifestlink, _ := os.Readlink("/tmp/cache_testTestFindManifest/cmdhash001/partial.1feb0822dbb0ef198ae7453b23513045bd677ef4b4d80e48b10983c9618fa598/manifest.111ac4231fc677d05a106ea0dac0ea5b628b1e8a8e0af287e266cf11f7efe94c")
			if manifestlink != "../manifest.403e4ce21cc1a8d2b887f3255db4559a9aa350034ac035ec90dd46e36c42e792" {
				os.RemoveAll("/tmp/cache_testTestFindManifest")
				t.Fail()
			}
		}else{
			os.RemoveAll("/tmp/cache_testTestFindManifest")
			t.Fail()
		}
		// test CopyOut from cache
		if utils.Exists("../test/hello.o") {
			os.Remove("../test/hello.o")
		}
		CopyOut(config, manifestFile)
		if !utils.Exists("../test/hello.o") {
			os.RemoveAll("/tmp/cache_testTestFindManifest")
			t.Fail()
		}
		// test find cache
	        manifestFile, _ = FindManifest(config, cmdhash)
		if manifestFile != "/tmp/cache_testTestFindManifest/cmdhash001/manifest.403e4ce21cc1a8d2b887f3255db4559a9aa350034ac035ec90dd46e36c42e792" {
			os.RemoveAll("/tmp/cache_testTestFindManifest")
			t.Fail()
		}
		// modify one inFile
		cmd = exec.Command("cp", "../test/sum.h", "../test/sum.h.bak")
		cmd.Run()
		ioutil.WriteFile("../test/sum.h", []byte("hello\n"), 0644)
	        manifestFile, _ = FindManifest(config, cmdhash)
		fmt.Println(manifestFile)
		cmd = exec.Command("cp", "../test/sum.h.bak", "../test/sum.h")
		cmd.Run()
		os.Remove("../test/sum.h.bak")
		// should no manifestFile found in cache
		if utils.Exists(manifestFile) {
			os.RemoveAll("/tmp/cache_testTestFindManifest")
			t.Fail()
		}
		// after revert back sum.bak, should find manifest in the cache 
		manifestFile, _ = FindManifest(config, cmdhash)
                if manifestFile != "/tmp/cache_testTestFindManifest/cmdhash001/manifest.403e4ce21cc1a8d2b887f3255db4559a9aa350034ac035ec90dd46e36c42e792" {
                        os.RemoveAll("/tmp/cache_testTestFindManifest")
                        t.Fail()
                }
        }else{
		os.RemoveAll("/tmp/cache_testTestFindManifest")
		t.Fail()
	}
	os.RemoveAll("/tmp/cache_testTestFindManifest")
}
