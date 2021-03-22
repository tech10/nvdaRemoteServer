package server

func Log(msg string, level int) {
	if level > loglevel {
		return
	}
	log_standard.Println(msg)
}

func Log_error(msg string) {
	log_error.Println(msg)
}
