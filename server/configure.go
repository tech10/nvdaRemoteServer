package server

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

var PS string = string(os.PathSeparator)

var confFile string
var DEFAULT_CONF_FILE string = ""
var DEFAULT_CONF_NAME string = "nvdaRemoteServer.json"
var DEFAULT_CONF_DIR string

var addresses AddressList
var DEFAULT_ADDRESS string = ":6837"

var cert string
var DEFAULT_CERT_FILE string = ""

var key string
var DEFAULT_KEY_FILE string = ""

var gencertfile string
var DEFAULT_GEN_CERT_FILE = ""

var logfile string
var DEFAULT_LOG_FILE string = ""

var loglevel int
var DEFAULT_LOG_LEVEL int = 0

const LOG_SILENT int = -1
const LOG_INFO int = 0
const LOG_CONNECTION int = 1
const LOG_CHANNEL int = 2
const LOG_DEBUG int = 3
const LOG_PROTOCOL int = 4

var motd string
var DEFAULT_MOTD string = ""

var motdAlwaysDisplay bool
var DEFAULT_MOTD_ALWAYS_DISPLAY bool = false

var sendOrigin bool
var DEFAULT_SEND_ORIGIN bool = true

var createDir bool
var DEFAULT_CREATE_DIR bool = false

var Launch bool
var DEFAULT_LAUNCH bool = true

var log_standard *log.Logger
var log_error *log.Logger

var Servers []*Server

var PID int
var PID_STR string
var pidfile string
var DEFAULT_PID_FILE string = ""

func init() {
	dcd, err := os.UserConfigDir()
	if err != nil {
		dcd = "."
	}
	DEFAULT_CONF_DIR = dcd + PS + "nvdaRemoteServer"
}

func Configure() error {
	PID = os.Getpid()
	PID_STR = strconv.Itoa(PID)

	flag.BoolVar(&createDir, "create", DEFAULT_CREATE_DIR, "Create directories upon any operation involving files being written to, or the working directory being changed.")

	flag.StringVar(&cert, "cert", DEFAULT_CERT_FILE, "SSL certificate file to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&key, "key", DEFAULT_KEY_FILE, "SSL key to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&gencertfile, "gen-cert-file", DEFAULT_GEN_CERT_FILE, "Generate a certificate file from the self-generated, self-signed SSL certificate. This file will only be created if you aren't loading your own certificate key files. The file will encode the key and certificate, packaging them both in a single .pem file.")

	flag.StringVar(&pidfile, "pid-file", DEFAULT_PID_FILE, "Create a PID file when the server has successfully started.")

	flag.Var(&addresses, "address", "Address the server will listen on in the format ip:port, such as \"0.0.0.0:6837\", \":6837\", \"[::]:6837\". The port must be between 1 and 65536. You can declare this parameter more than once for multiple listen addresses.")

	flag.IntVar(&loglevel, "log-level", DEFAULT_LOG_LEVEL, "Choose what log level you wish to use. Any value below -1 will be ignored.")
	flag.StringVar(&logfile, "log-file", DEFAULT_LOG_FILE, "Choose what log file you wish to use in addition to logging output to the console. If the file can't be created or open for writing, the program will fall back to console logging only.")

	flag.StringVar(&motd, "motd", DEFAULT_MOTD, "Display a message of the day for the server.")
	flag.BoolVar(&motdAlwaysDisplay, "motd-always-display", DEFAULT_MOTD_ALWAYS_DISPLAY, "Force the message of the day to be displayed upon each connection to the server, even if it hasn't changed.")

	flag.BoolVar(&sendOrigin, "send-origin", DEFAULT_SEND_ORIGIN, "Send an origin message from every message received by a client.")

	flag.BoolVar(&Launch, "launch", DEFAULT_LAUNCH, "Launch the server.")

	flag.Parse()

	log_init(logfile)

	Log(LOG_INFO, "Initializing configuration.")

	generate := false
	var config *tls.Config
	var err error

	if cert != "" && !fileExists(cert) {
		Log(LOG_INFO, "The certificate file at "+cert+" does not exist.")
		generate = true
	}
	if key != "" && !fileExists(key) {
		Log(LOG_INFO, "The key file at "+key+" does not exist.")
		generate = true
	}
	if cert == "" || key == "" {
		generate = true
	}

	if generate {
		Log(LOG_DEBUG, "Attempting to generate self-signed SSL certificate.")
		config, err = gen_cert()
		if err != nil {
			Log_error("Unable to generate self-signed certificate.\r\n" + err.Error() + "\r\nUnable to start server.")
			Launch_fail()
			return err
		}
		Log(LOG_DEBUG, "SSL certificate generated.")
	} else {
		if gencertfile != "" {
			Log(LOG_INFO, "The server has not generated its own self-signed certificate, and the -gen-certfile parameter is set to "+gencertfile+". This parameter will be ignored.")
		}
		cert, cerr := tls.LoadX509KeyPair(cert, key)
		if cerr != nil {
			Log_error("Error loading certificate and key files.\r\n" + cerr.Error() + "\r\nUnable to start server.")
			Launch_fail()
			return cerr
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	config.PreferServerCipherSuites = false
	config.MaxVersion = tls.VersionTLS12
	config.InsecureSkipVerify = true

	if motd != "" {
		logstr := "The server will display the following message of the day:\r\n" + motd
		if motdAlwaysDisplay {
			logstr += "\r\nThe server will tell each client to display this message of the day upon each connection."
		}
		Log(LOG_DEBUG, logstr)
	}

	if motd == "" && motdAlwaysDisplay {
		Log(LOG_INFO, "The server has been told to always display a message of the day, but no message of the day has been set. The -motd-always-display parameter will be reset to false.")
		motdAlwaysDisplay = false
	}

	if !sendOrigin {
		Log(LOG_INFO, "The server is configured to send no origin message to other clients, which may improve performance slightly, but impact the useability of your server when the origin field is required.")
	}

	if !Launch {
		Log(LOG_INFO, "The server will not be launched. Shutting down.")
		return errors.New("Server launch parameter set to false.")
	}

	if len(addresses) == 0 {
		addresses = make(AddressList, 1)
		addresses[0] = DEFAULT_ADDRESS
	}

	Servers = make([]*Server, len(addresses))
	for i, addr := range addresses {
		Servers[i] = NewWithTLSConfig(addr, config)
		Log(LOG_DEBUG, "Starting server listening on address "+addr)
	}

	return nil
}

func Start() int {
	var num int = 0
	var err error

	for i := range Servers {
		err = Servers[i].Listen()
		if err != nil {
			Log_error("Unable to listen on address " + Servers[i].address + ".\r\n" + err.Error())
			Servers[i] = nil
			continue
		}
		num++
	}
	if num == 0 {
		Servers = nil
		return num
	}

	Log(LOG_DEBUG, "Number of servers started: "+strconv.Itoa(num))
	go signals_init()
	return num
}

func Launch_fail() {
	if !Launch {
		os.Exit(1)
	}
}

func cfg_default() *Cfg {
	return &Cfg{
		PidFile:           DEFAULT_PID_FILE,
		LogFile:           DEFAULT_LOG_FILE,
		LogLevel:          DEFAULT_LOG_LEVEL,
		Addresses:         AddressList{DEFAULT_ADDRESS},
		Cert:              DEFAULT_CERT_FILE,
		Key:               DEFAULT_KEY_FILE,
		Motd:              DEFAULT_MOTD,
		MotdAlwaysDisplay: DEFAULT_MOTD_ALWAYS_DISPLAY,
		SendOrigin:        DEFAULT_SEND_ORIGIN,
	}
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

func cfg_is_default(c *Cfg) bool {
	if !default_pid_file(c.PidFile) {
		return false
	}
	if !default_log_file(c.LogFile) {
		return false
	}
	if !default_log_level(c.LogLevel) {
		return false
	}
	if !default_addresses(c.Addresses) {
		return false
	}
	if !default_cert_file(c.Cert) {
		return false
	}
	if !default_key_file(c.Key) {
		return false
	}
	if !default_motd(c.Motd) {
		return false
	}
	if !default_motd_always_display(c.MotdAlwaysDisplay) {
		return false
	}
	if !default_send_origin(c.SendOrigin) {
		return false
	}
	return true
}

func conf_file_check() ([]byte, error) {
	if !default_conf_file(confFile) {
		return file_read(confFile)
	}
	d, err := file_read(DEFAULT_CONF_NAME)
	if err != nil {
		return file_read(DEFAULT_CONF_DIR + PS + DEFAULT_CONF_NAME)
	}
	return d, err
}

func conf_file_error(err error) error {
	if confFile == DEFAULT_CONF_FILE {
		return nil
	}
	return errors.New("Unable to read configuration file.\n" + err.Error())
}

func conf_file() error {
	d, err := conf_file_check()
	if err != nil {
		return conf_file_error(err)
	}
	cfg := cfg_default()
	err = cfg_read(d, cfg)
	if err != nil {
		return conf_file_error(err)
	}
	if cfg_is_default(cfg) {
		return nil
	}
	conf_set(cfg)
	return nil
}

func conf_set(c *Cfg) {

}
