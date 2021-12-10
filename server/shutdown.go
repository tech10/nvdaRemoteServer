package server

import (
	"os"

	"github.com/tech10/panichandler"
)

func Shutdown() {
	PidfileClear()
}

var PanicHandle panichandler.Capture = panichandler.Capture{
	F: func(i *panichandler.Info) {
		Log_error("PANIC\n", i.String())
		Shutdown()
		os.Exit(2)
	},
	ExitCode: panichandler.ExitCode,
}
