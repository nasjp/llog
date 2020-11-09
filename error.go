package main

import (
	"fmt"
)

type errInvalidArg struct{ Arg []string }

func (e errInvalidArg) Error() string {
	return fmt.Sprintf("invalid args: %v\n%s", e.Arg, helpMessage)
}

type errInvalidCmd struct{ Cmd string }

func (e errInvalidCmd) Error() string {
	return fmt.Sprintf("invalid command: %s\n%s", e.Cmd, helpMessage)
}

type errInternal struct{ Err error }

func (e errInternal) Error() string {
	return fmt.Sprintf("internal llog error: %s", e.Err.Error())
}
