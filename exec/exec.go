package exec

import (
	"bufio"
	"config"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"utils"
)

func ParseWiskTrackFile(trackfile string) (infiles []string, outfiles []string) {
	file, err := os.Open(trackfile)
	if err != nil {
		return infiles, outfiles
	}
	defer file.Close()

	context := map[string]string{}
	var jsondata []interface{}
	var line string
	var parts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		parts = strings.SplitN(line, " ", 3)
		if parts[1] == "READS" {
			json.Unmarshal([]byte(parts[2]), &jsondata)
			// fmt.Println("READS: ", jsondata)
			if opfile, ok := jsondata[0].(string); ok {
				if !filepath.IsAbs(opfile) {
					opfile = filepath.Join(context[parts[0]], opfile)
				} else {
					opfile = filepath.Join(opfile, "")
				}
				infiles = append(infiles, opfile)
			} else {
				panic(ok)
			}
		} else if parts[1] == "WRITES" {
			json.Unmarshal([]byte(parts[2]), &jsondata)
			// fmt.Println("WRITES: ", jsondata)
			if opfile, ok := jsondata[0].(string); ok {
				if strings.HasPrefix(opfile, "/dev/") {
					continue
				}
				if !filepath.IsAbs(opfile) {
                                        // opfile is based on WSROOT
					opfile = filepath.Join(context[parts[0]], opfile)
				} else {
					opfile = filepath.Join(opfile, "")
				}
				outfiles = append(outfiles, opfile)
			} else {
				panic(ok)
			}
		} else if parts[1] == "CALLS" {
			json.Unmarshal([]byte(parts[2]), &jsondata)
			if oplist, ok := jsondata[0].([]interface{}); ok {
				if uuid, ok := oplist[1].(string); ok {
					// fmt.Println("UUID: ", uuid)
					// if oplist, ok := jsondata[2].([]interface{}); ok {
                                        // 2 is CWD, 3 is WSROOT
					if oplist, ok := jsondata[3].([]interface{}); ok {
						if cwd, ok := oplist[1].(string); ok {
							// fmt.Println("CWD: ", cwd)
							// fmt.Println("WSROOT: ", cwd)
							context[uuid] = cwd
						} else {
							panic(ok)
						}
					} else {
						panic(ok)
					}
				} else {
					panic(ok)
				}
			} else {
				panic(ok)
			}
		}
	}
	// fmt.Println("Infiles: ", infiles)
	// fmt.Println("Outfiles: ", outfiles)
	return infiles, outfiles
}

func RunCmd(conf config.Config, cmdhash string, cmd []string) (exitcode int, logfile string, infiles []string, outfiles []string) {
	fmt.Println("Executing: ", cmd)
	fmt.Println("Hash: ", cmdhash)
	logfile = fmt.Sprintf("/tmp/wisktrack/wiskcachecmdrun.%s.log", cmdhash)
	trackfile := fmt.Sprintf("/tmp/wisktrack/wisktrack.%s.file", cmdhash)
	os.Remove(trackfile)
	if !utils.Exists(filepath.Dir(logfile)) {
		os.MkdirAll(filepath.Dir(logfile), 0775)
	}

	out, err := os.Create(logfile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = out
	command.Stderr = out
	command.Env = append(
		os.Environ(),
		"LD_PRELOAD=/ws/sarvi-sjc/wisktrack/${LIB}/libwisktrack.so",
		"WISK_CONFIG=",
		"WISK_TRACE=%s/wisktrace.log",
		fmt.Sprintf("WISK_TRACK=/tmp/wisktrack/wisktrack.%s.file", cmdhash),
		fmt.Sprintf("WISK_WSROOT=%s", conf.BaseDir),
		fmt.Sprintf("WISK_TRACE=%s/wisktrace.log", conf.BaseDir),
	)
	// fmt.Println(command.Env, command)
	// if err := command.Run(); err != nil {
	// 	log.Fatal(err)
	// 	exitcode = 1
	// 	return
	// }
	// fmt.Println("Run: ", out.String())
	fmt.Println("Run Trackfile: ", trackfile)
	err = command.Start()
	if err != nil {
		panic(err)
	}
	command.Wait()
	exitcode = 0
	infiles, outfiles = ParseWiskTrackFile(trackfile)
	fmt.Println("Run Logfile: ", logfile)
	fmt.Println("Run Infiles: ", infiles)
	fmt.Println("Run Outfile: ", outfiles)
	return
}
