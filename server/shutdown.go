package server

import (
	"github.com/tech10/panic_handler"
	"os"
)

func Shutdown() {
	PidfileClear()
}

var PanicHandle panic_handler.HandlerFunc = func(i *panic_handler.Info) {
	Log_error("PANIC\n", i.String())
	Shutdown()
	os.Exit(2)
}
