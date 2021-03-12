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
			cmd: []string{"/bin/bash", "-c", "touch /tmp/testwiskcache.input; cat /tmp/testwiskcache.input > /tmp/testwiskcache.output"},
		},
	}
	for _, tc := range testCases {
		fmt.Println("\tSubTest: ", tc.d)
		exitcode, infiles, outfiles := RunCmd(conf, "asdasdasd", tc.cmd)
		fmt.Println("Exec Failed: ", exitcode, infiles, outfiles, reflect.DeepEqual(infiles, []string{}))
		// if exitcode != 0 || !reflect.DeepEqual(infiles, []string{}) || !reflect.DeepEqual(outfiles, []string{}) {
		if exitcode != 0 { // || !reflect.DeepEqual(infiles, []string{}) || !reflect.DeepEqual(outfiles, []string{}) {
			fmt.Println("Exec Failed: ", exitcode, tc.cmd, infiles, outfiles)
			t.Fail()
		}

	}
}
