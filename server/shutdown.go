package server

import (
	"os"
	"runtime/debug"
)

func Shutdown() {
	PidfileClear()
}

func PanicHandle() {
	r := recover()
	if r == nil {
		return
	}
	Log_error("PANIC\n", r, "\n", string(debug.Stack()))
	Shutdown()
	os.Exit(2)
}
