package server

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

const DEFAULT_ADDRESS string = ":6837"

var addresses AddressList

var Cert string

const DEFAULT_CERT_FILE string = ""

var Key string

const DEFAULT_KEY_FILE string = ""

var gencertfile string

const DEFAULT_GEN_CERT_FILE = ""

var logfile string

const DEFAULT_LOG_FILE string = ""

var loglevel int

const DEFAULT_LOGLEVEL int = 0
const LOG_SILENT int = -1
const LOG_INFO int = 0
const LOG_CONNECTION int = 1
const LOG_CHANNEL int = 2
const LOG_DEBUG int = 3
const LOG_PROTOCOL int = 4

var motd string

const DEFAULT_MOTD string = ""

var motdAlwaysDisplay bool

const DEFAULT_MOTD_ALWAYS_DISPLAY bool = false

var sendOrigin bool

const DEFAULT_SEND_ORIGIN bool = true

var createDir bool

const DEFAULT_CREATE_DIR bool = false

var Launch bool

const DEFAULT_LAUNCH bool = true

var log_standard *log.Logger
var log_error *log.Logger

var Servers []*Server

var PID int
var PID_STR string
var pidfile string

const DEFAULT_PID_FILE string = ""

func Configure() error {
	PID = os.Getpid()
	PID_STR = strconv.Itoa(PID)

	flag.BoolVar(&createDir, "create", DEFAULT_CREATE_DIR, "Create directories upon any operation involving files being written to, or the working directory being changed.")

	flag.StringVar(&Cert, "cert", DEFAULT_CERT_FILE, "SSL certificate file to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&Key, "key", DEFAULT_KEY_FILE, "SSL key to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&gencertfile, "gen-cert-file", DEFAULT_GEN_CERT_FILE, "Generate a certificate file from the self-generated, self-signed SSL certificate. This file will only be created if you aren't loading your own certificate key files. The file will encode the key and certificate, packaging them both in a single .pem file.")

	flag.StringVar(&pidfile, "pid-file", DEFAULT_PID_FILE, "Create a PID file when the server has successfully started.")

	flag.Var(&addresses, "address", "Address the server will listen on in the format ip:port, such as \"0.0.0.0:6837\", \":6837\", \"[::]:6837\". The port must be between 1 and 65536. You can declare this parameter more than once for multiple listen addresses.")

	flag.IntVar(&loglevel, "log-level", DEFAULT_LOGLEVEL, "Choose what log level you wish to use. Any value below -1 will be ignored.")
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

	if Cert != "" && !fileExists(Cert) {
		Log(LOG_INFO, "The certificate file at "+Cert+" does not exist.")
		generate = true
	}
	if Key != "" && !fileExists(Key) {
		Log(LOG_INFO, "The key file at "+Key+" does not exist.")
		generate = true
	}
	if Cert == "" || Key == "" {
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
		cert, cerr := tls.LoadX509KeyPair(Cert, Key)
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
