package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	var dir string

	pflag.StringVarP(&dir, "dir", "d", "", "target directory")
	pflag.Parse()

	if len(pflag.Args()) != 1 {
		return errInvalidArg{Arg: pflag.Args()}
	}

	switch cmd := pflag.Args()[0]; cmd {
	case cmdGen:
		return gen(dir)
	case cmdHelp:
		fmt.Println(helpMessage)
	default:
		return errInvalidCmd{Cmd: cmd}
	}

	return nil
}
