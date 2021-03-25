package server

func Log(level int, msg ...interface{}) {
	if level > loglevel {
		return
	}
	log_standard.Println(msg...)
}

func Log_error(msg ...interface{}) {
	log_error.Println(msg...)
}
