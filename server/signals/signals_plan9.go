// +build plan9

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait() chan os.Signal {
	//Signal notifiers
	kill := make(chan os.Signal, 2)
	signal.Notify(kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.Note("quit"))
	return kill
}
