// nvdaRemoteServer

package main

import (
	. "github.com/tech10/nvdaRemoteServer/server"
	"os"
	"strconv"
	"sync"
)

func main() {
	err := Configure()
	if err != nil {
		return
	}
	num := Start()
	if num == 0 {
		return
	}
	Log("Server started. Running under PID "+strconv.Itoa(os.Getpid()), LOG_INFO)
	wait()
	Log("Server shutdown complete.", LOG_INFO)
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
