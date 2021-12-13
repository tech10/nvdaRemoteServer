package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	. "github.com/tech10/nvdaRemoteServer/server"
)

var Version string = "development"

func main() {
	Version = strings.TrimPrefix(versionSetter(), "v")
	args()

	defer Log_close()
	err := Configure()
	if err != nil {
		if Launch {
			Log_close()
			os.Exit(1)
		}
		return
	}
	num := Start()
	if num == 0 {
		Log_error("No servers started. Shutting down.")
		Log_close()
		os.Exit(1)
	}
	defer PanicHandle.Catch()
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

func versionSetter() string {
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}
	m := i.Main
	if m.Sum != "" {
		return m.Version
	}
	return Version
}
