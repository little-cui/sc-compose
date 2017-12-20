package command

import (
	"fmt"
	"os"
)

const (
	ErrExit = iota
	ErrInvalidParam
	ErrUnmarshal
)

func ExitWithError(err int, msgs ...interface{}) {
	fmt.Fprintln(os.Stderr, "error:", err, fmt.Sprint(msgs...))
	os.Exit(err)
}
