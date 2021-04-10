package server

import (
	"log"
	"os"
)

func Log(level int, msg ...interface{}) {
	if level > loglevel {
		return
	}
	log_standard.Println(msg...)
}

func Log_error(msg ...interface{}) {
	log_error.Println(msg...)
}

func log_init() {
	log_standard = log.New(os.Stdout, "", log.LstdFlags)
	log_error = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags)
}
