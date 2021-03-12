package exec

import (
	"bufio"
	"config"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func jsonlist(json string) string {
	json = strings.Trim(json, "[] ")
	// fmt.Println("Trimed: ", json)
	return strings.Split(json, ",")[0]
}

func ParseWiskTrackFile(trackfile string) (infiles []string, outfiles []string) {
	file, err := os.Open(trackfile)
	if err != nil {
		return infiles, outfiles
	}
	defer file.Close()

	var line string
	var parts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		parts = strings.SplitN(line, " ", 3)
		if parts[1] == "READS" {
			// fmt.Println("OP: ", parts[1], parts[2])
			infiles = append(infiles, jsonlist(parts[2]))
		} else if parts[1] == "WRITES" {
			// fmt.Println("OP: ", parts[1], parts[2])
			outfiles = append(outfiles, jsonlist(parts[2]))
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
	fmt.Println("Run Outfiles: ", infiles)
	fmt.Println("Run Trackfile: ", outfiles)
	return
}
