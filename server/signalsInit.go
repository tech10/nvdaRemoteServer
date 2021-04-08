// +build !plan9

package server

import (
	"os"
	"os/signal"
	"syscall"
)

func signalsWait() chan os.Signal {
	//Signal notifiers
	kill := make(chan os.Signal, 2)
	signal.Notify(kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return kill
}
