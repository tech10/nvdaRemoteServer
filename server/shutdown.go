package server

import (
	"os"

	"github.com/tech10/panic_handler"
)

func Shutdown() {
	PidfileClear()
}

var PanicHandle panic_handler.Capture = panic_handler.Capture{
	F: func(i *panic_handler.Info) {
		Log_error("PANIC\n", i.String())
		Shutdown()
		os.Exit(2)
	},
	ExitCode: panic_handler.ExitCode,
}
