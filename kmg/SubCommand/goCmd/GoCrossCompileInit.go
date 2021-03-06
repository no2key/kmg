package goCmd

import (
	"fmt"
	"path/filepath"

	"os/exec"
	"runtime"
	"strings"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
)

func runGoCrossCompileInit() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	GOROOT := kmgc.GOROOT
	if GOROOT == "" {
		//guess GOROOT
		out, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
		kmgConsole.ExitOnErr(err)
		GOROOT = strings.TrimSpace(string(out))
		if GOROOT == "" {
			kmgConsole.ExitOnErr(fmt.Errorf("you must set $GOROOT in environment to use GoCrossComplieInit"))
		}
	}
	var makeShellArgs []string
	var makeShellName string
	runCmdPath := filepath.Join(GOROOT, "src")
	if runtime.GOOS == "windows" {
		makeShellName = "cmd"
		makeShellArgs = []string{"/C", filepath.Join(GOROOT, "src", "make.bat"), "--no-clean"}
	} else {
		makeShellName = filepath.Join(GOROOT, "src", "make.bash")
		makeShellArgs = []string{"--no-clean"}
	}
	for _, target := range kmgc.CrossCompileTarget {
		err = kmgCmd.CmdSlice(append([]string{makeShellName}, makeShellArgs...)).
			MustSetEnv("GOOS", target.GetGOOS()).
			MustSetEnv("GOARCH", target.GetGOARCH()).
			SetDir(runCmdPath).
			Run()
		kmgConsole.ExitOnErr(err)
	}
	return
}
