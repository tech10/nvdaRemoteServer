package server

import "github.com/tech10/nvdaRemoteServer/signals"

func signals_init() {
	<-signals.Wait()
	Log(LOG_INFO, "Signal received to shut down.")
	StopServers()
}
