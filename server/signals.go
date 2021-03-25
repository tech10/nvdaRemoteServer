package server

import (
	"os"
	"os/signal"
	"syscall"
)

func signals_init() {
	//Signal notifiers
	kill := make(chan os.Signal, 2)
	signal.Notify(kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-kill
	Log(LOG_INFO, "Signal received to shut down.")
	StopServers()
}
