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
	Log("Signal received to shut down.", LOG_INFO)
	StopServers()
}
