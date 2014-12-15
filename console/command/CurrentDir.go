package command

import (
	"github.com/bronze1man/kmg/kmgConsole"
	"os"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "CurrentDir",
		Desc:   "get current dir(usefull in cygwin)",
		Runner: runCurrentDir,
	})
}

func runCurrentDir() {
	wd, err := os.Getwd()
	exitOnErr(err)
	_, err = os.Stdout.Write([]byte(wd))
	exitOnErr(err)
}
