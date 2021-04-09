package server

import "github.com/tech10/server/signals"

func signals_init() {
	<-signals.Wait()
	Log(LOG_INFO, "Signal received to shut down.")
	StopServers()
}
