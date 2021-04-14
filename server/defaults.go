package server

import (
	"os"
)

var PS string = string(os.PathSeparator)

var DEFAULT_CONF_FILE string = ""
var DEFAULT_CONF_NAME string = "nvdaRemoteServer.json"
var DEFAULT_CONF_DIR string

var DEFAULT_ADDRESS string = ":6837"

var DEFAULT_CERT_FILE string = ""
var DEFAULT_KEY_FILE string = ""
var DEFAULT_GEN_CERT_FILE = ""

var DEFAULT_LOG_FILE string = ""

var DEFAULT_LOG_LEVEL int = 0

const LOG_SILENT int = -1
const LOG_INFO int = 0
const LOG_CONNECTION int = 1
const LOG_CHANNEL int = 2
const LOG_DEBUG int = 3
const LOG_PROTOCOL int = 4

var DEFAULT_MOTD string = ""
var DEFAULT_MOTD_ALWAYS_DISPLAY bool = false

var DEFAULT_SEND_ORIGIN bool = true

var DEFAULT_CREATE_DIR bool = false

var DEFAULT_LAUNCH bool = true

var DEFAULT_PID_FILE string = ""

func init() {
	dcd, err := os.UserConfigDir()
	if err != nil {
		dcd = "."
	}
	DEFAULT_CONF_DIR = dcd + PS + "nvdaRemoteServer"
}

func default_conf_file(p string) bool {
	return (p == DEFAULT_CONF_FILE)
}

func default_pid_file(p string) bool {
	return (p == DEFAULT_PID_FILE)
}

func default_log_file(p string) bool {
	return (p == DEFAULT_LOG_FILE)
}

func default_log_level(p int) bool {
	return (p == DEFAULT_LOG_LEVEL)
}

func default_addresses(p AddressList) bool {
	if len(p) == 1 && p[0] == DEFAULT_ADDRESS {
		return true
	}
	if len(p) == 0 {
		return true
	}
	return false
}

func default_cert_file(p string) bool {
	return (p == DEFAULT_CERT_FILE)
}

func default_key_file(p string) bool {
	return (p == DEFAULT_KEY_FILE)
}

func default_motd(p string) bool {
	return (p == DEFAULT_MOTD)
}

func default_motd_always_display(p bool) bool {
	return (p == DEFAULT_MOTD_ALWAYS_DISPLAY)
}

func default_send_origin(p bool) bool {
	return (p == DEFAULT_SEND_ORIGIN)
}
