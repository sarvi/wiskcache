package cache

import (
	"fmt"
	"testing"
	"sync"
	"os/exec"
	"strings"
	"os"
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
