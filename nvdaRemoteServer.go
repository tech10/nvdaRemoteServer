package main

import (
	. "github.com/tech10/nvdaRemoteServer/server"
	"os"
	"strconv"
	"strings"
	"sync"
)

var Version string = "development"

func main() {
	// Log panics
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		Log_error(r)
	}()

	Version = strings.TrimPrefix(Version, "v")
	err := Configure()
	if err != nil {
		if Launch {
			os.Exit(1)
		}
		return
	}
	num := Start()
	if num == 0 {
		Log_error("No servers started. Shutting down.")
		os.Exit(1)
	}
	Log(LOG_INFO, "Server started. Running under PID "+strconv.Itoa(os.Getpid())+". Server version "+Version)
	wait()
	Log(LOG_INFO, "Server shutdown complete.")
}

func wait() {
	var wg sync.WaitGroup
	for _, s := range Servers {
		if s == nil {
			continue
		}
		wg.Add(1)
		go func(sv *Server) {
			sv.Wait()
			wg.Done()
		}(s)
	}
	wg.Wait()
}
