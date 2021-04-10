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
	Version = strings.TrimPrefix(Version, "v")
	err := Configure()
	if err != nil {
		os.Exit(1)
	}
	num := Start()
	if num == 0 {
		os.Exit(1)
	}
	Log(LOG_INFO, "Server started. Running under PID "+strconv.Itoa(os.Getpid())+". Server version "+Version)
	wait()
	Log(LOG_INFO, "Server shutdown complete.")
}

func wait() {
	var wg sync.WaitGroup
	if S4 != nil {
		wg.Add(1)
		go func() {
			S4.Wait()
			wg.Done()
		}()
	}
	if S6 != nil {
		wg.Add(1)
		go func() {
			S6.Wait()
			wg.Done()
		}()
	}
	if SAll != nil {
		wg.Add(1)
		go func() {
			SAll.Wait()
			wg.Done()
		}()
	}
	wg.Wait()
}
