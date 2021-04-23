package main

import (
	"fmt"
	. "github.com/tech10/nvdaRemoteServer/server"
	"os"
	"strings"
	"sync"
)

var Version string = "development"

func main() {
	Version = strings.TrimPrefix(Version, "v")
	args()
	// Log panics
	defer PanicHandle()

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
	PidfileSet()
	Log(LOG_INFO, "Server started. Running under PID "+PID_STR+". Server version "+Version)
	wait()
	Shutdown()
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

func args() {
	if len(os.Args) < 2 {
		return
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(Version)
		os.Exit(0)
	default:
		return
	}
}
