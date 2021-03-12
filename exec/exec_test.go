package exec

import (
	"config"
	"fmt"
	"reflect"
	"testing"
)

func TestRunCmd(t *testing.T) {
	conf := config.Config{BaseDir: "/ws/sarvi-sjc/wiskcache/exec"}
	testCases := []struct {
		d   string
		cmd []string
	}{
		{
			d:   "relative dirs",
			cmd: []string{"/bin/bash", "-c", "echo \"Hello World\" ; cat tests/file1.in > tests/file.out ; cat tests/file2.in >> tests/file.out ; cat tests/file2.in >> tests/file.out"},
		},
	}
	for _, tc := range testCases {
		fmt.Println("\tSubTest: ", tc.d)
		exitcode, logfile, infiles, outfiles := RunCmd(conf, "asdasdasd", tc.cmd)
		fmt.Println("Exec Failed: ", exitcode, logfile, infiles, outfiles, reflect.DeepEqual(infiles, []string{}))
		// if exitcode != 0 || !reflect.DeepEqual(infiles, []string{}) || !reflect.DeepEqual(outfiles, []string{}) {
		if exitcode != 0 { // || !reflect.DeepEqual(infiles, []string{}) || !reflect.DeepEqual(outfiles, []string{}) {
			fmt.Println("Exec Failed: ", exitcode, tc.cmd, infiles, outfiles)
			t.Fail()
		}

	}
}
