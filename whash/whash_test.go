package whash

import (
	"config"
	"fmt"
	"testing"
)

type TestData struct {
	conf config.Config
	env  map[string]string
	cmd  []string
}

type TestDataPair struct {
	d string
	l TestData
	r TestData
}

func TestHashMustMatch(t *testing.T) {
	testCases := []TestDataPair{
		{
			d: "Same all",
			l: TestData{
				conf: config.Config{
					BaseDir: "/a/b",
					Envars:  []string{"CWD"},
					Tools:   []config.Tool{},
					ToolIdx: -1,
				},
				env: map[string]string{"CWD": "/a/b/c"},
				cmd: []string{"g++", "-o", "file.o", "file.c"},
			},
			r: TestData{
				conf: config.Config{
					BaseDir: "/a/b",
					Envars:  []string{"CWD"},
					Tools:   []config.Tool{},
					ToolIdx: -1,
				},
				env: map[string]string{"CWD": "/a/b/c"},
				cmd: []string{"g++", "-o", "file.o", "file.c"},
			},
		},
		// {
		// 	d: "Same command and env, absolute vs relative path",
		// 	l: TestData{
		// 		conf: config.Config{
		// 			BaseDir: "/a/b",
		// 			Envars:  []string{"CWD"},
		// 			Tools:   []config.Tool{},
		// 			ToolIdx: -1,
		// 		},
		// 		env: map[string]string{"CWD": "/a/b/c"},
		// 		cmd: []string{"g++", "-o", "file.o", "/a/b/c/file.c"},
		// 	},
		// 	r: TestData{
		// 		conf: config.Config{
		// 			BaseDir: "/a/b",
		// 			Envars:  []string{"CWD"},
		// 			Tools:   []config.Tool{},
		// 			ToolIdx: -1,
		// 		},
		// 		env: map[string]string{"CWD": "/a/b/c"},
		// 		cmd: []string{"g++", "-o", "file.o", "file.c"},
		// 	},
		// },
	}
	for _, tc := range testCases {
		fmt.Println("\tSubTest: ", tc.d)
		h1, e1 := cmdhash(tc.l.conf, tc.l.env, tc.l.cmd)
		h2, e2 := cmdhash(tc.r.conf, tc.r.env, tc.r.cmd)
		// fmt.Println(h1, h2)
		if h1 != h2 || e1 != nil || e2 != nil {
			t.Fail()
		}
	}
}

func TestHashMustNotMatch(t *testing.T) {
	testCases := []TestDataPair{
		{
			d: "Same command, different CWD",
			l: TestData{
				conf: config.Config{
					BaseDir: "/a/b",
					Envars:  []string{"CWD"},
					Tools:   []config.Tool{},
					ToolIdx: -1,
				},
				env: map[string]string{"CWD": "/a/b/c"},
				cmd: []string{"g++", "-o", "file.o", "file.c"},
			},
			r: TestData{
				conf: config.Config{
					BaseDir: "/a/b",
					Envars:  []string{"CWD"},
					Tools:   []config.Tool{},
					ToolIdx: -1,
				},
				env: map[string]string{"CWD": "/a/b/d"},
				cmd: []string{"g++", "-o", "file.o", "file.c"},
			},
		},
	}
	for _, tc := range testCases {
		fmt.Println("\tSubTest: ", tc.d)
		h1, e1 := cmdhash(tc.l.conf, tc.l.env, tc.l.cmd)
		h2, e2 := cmdhash(tc.r.conf, tc.r.env, tc.r.cmd)
		// fmt.Println(h1, h2)
		if h1 == h2 || e1 != nil || e2 != nil {
			t.Fail()
		}
	}
}
