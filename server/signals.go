package server

import "github.com/tech10/nvdaRemoteServer/signals"

func signals_init() {
	sig := <-signals.Wait()
	Log(LOG_INFO, "Signal received to shut down. Received signal "+sig.String())
	StopServers()
}
