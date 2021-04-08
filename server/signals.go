package server

func signals_init() {
	<-signalsWait()
	Log(LOG_INFO, "Signal received to shut down.")
	StopServers()
}
