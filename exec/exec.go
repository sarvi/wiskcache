package exec

import (
	"config"
	"fmt"
	"os"
	"os/exec"
)

func RunCmd(conf config.Config, cmdhash string, cmd []string) (exitcode int, infiles []string, outfiles []string) {
	outfile := fmt.Sprintf("/tmp/wisktrack/wiskcachecmdrun.%s.log", cmdhash)
	trackfile := fmt.Sprintf("/tmp/wisktrack/wisktrack.%s.file", cmdhash)
	os.Remove(trackfile)
	out, err := os.Create(outfile)
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
	fmt.Println("Run Outfile: ", outfile)
	fmt.Println("Run Trackfile: ", trackfile)
	err = command.Start()
	if err != nil {
		panic(err)
	}
	command.Wait()
	exitcode = 0
	return
}
