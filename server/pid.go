package server

import (
	"os"
)

func PidfileSet() {
	if pidfile == "" {
		return
	}
	Log(LOG_DEBUG, "Attempting to write process ID to PID file "+pidfile)
	err := file_rewrite(pidfile, []byte(PID_STR))
	if err != nil {
		Log(LOG_DEBUG, "Failed to write PID file.\n"+err.Error())
		pidfile = ""
	}
	Log(LOG_DEBUG, "Successfully wrote PID file.")
}

func PidfileClear() {
	if pidfile == "" {
		return
	}
	Log(LOG_DEBUG, "Removing PID file "+pidfile)
	err := os.Remove(pidfile)
	if err != nil {
		Log(LOG_DEBUG, "Failed to remove "+pidfile+"\n"+err.Error())
		pidfile = ""
		return
	}
	pidfile = ""
	Log(LOG_DEBUG, "Successfully removed PID file.")
}
