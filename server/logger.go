package server

import (
	"io"
	"log"
	"os"
	"sync"
)

var ll sync.Mutex

var (
	log_standard *log.Logger
	log_error    *log.Logger
	log_file     *os.File
)

func Log(level int, msg ...interface{}) {
	if level > loglevel {
		return
	}
	ll.Lock()
	defer ll.Unlock()
	log_standard.Println(msg...)
}

func Log_error(msg ...interface{}) {
	ll.Lock()
	defer ll.Unlock()
	log_error.Println(msg...)
}

func log_init(file string) {
	if file == "" {
		log_standard = log.New(os.Stdout, "", log.LstdFlags)
		log_error = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags)
		return
	}
	file = fullPath(file)
	if !Launch {
		log_init("")
		Log(LOG_INFO, "The log file at "+file+" will not be written, as the server has been told not to launch.")
		return
	}
	w, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log_init("")
		Log_error("Unable to open log file " + file + " for writing.\r\n" + err.Error())
		return
	}
	log_standard = log.New(io.MultiWriter(os.Stdout, w), "", log.LstdFlags)
	log_error = log.New(io.MultiWriter(os.Stderr, w), "[ERROR]: ", log.LstdFlags)
	log_file = w
}

func Log_close() {
	ll.Lock()
	defer ll.Unlock()
	if log_file == nil {
		return
	}
	_ = log_file.Sync()
	_ = log_file.Close()
	log_file = nil
}
