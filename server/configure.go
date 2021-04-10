package server

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

var ip4 string
var ip6 string
var port int

const DEFAULT_PORT int = 6837

var Cert string
var Key string
var gencertfile string

var logfile string

var loglevel int

const DEFAULT_LOGLEVEL int = 0
const LOG_SILENT int = -1
const LOG_INFO int = 0
const LOG_CONNECTION int = 1
const LOG_CHANNEL int = 2
const LOG_DEBUG int = 3
const LOG_PROTOCOL int = 4

var motd string
var motdForceDisplay bool

var sendOrigin bool

var Launch bool

var log_standard *log.Logger
var log_error *log.Logger

var S4 *Server
var S6 *Server
var SAll *Server

func Configure() error {
	flag.StringVar(&Cert, "cert", "", "SSL certificate file to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&Key, "key", "", "SSL key to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&gencertfile, "gen-certfile", "", "Generate a certificate file from the self-generated, self-signed SSL certificate. This file will only be created if you aren't loading your own certificate key files. The file will encode the key and certificate, packaging them both in a single .pem file.")

	flag.StringVar(&ip4, "ip4", "", "IPV4 address for the server to listen for connections on. This can be blank if desired, in which case, the server will listen on all IPV4 addresses.")
	flag.StringVar(&ip6, "ip6", "", "IPV6 address for the server to listen for connections on. This can be blank if desired, in which case, the server will listen on all IPV6 addresses.")
	flag.IntVar(&port, "port", DEFAULT_PORT, "The port that the server will listen for connections on. This can be blank if desired, in which case, the server will listen for connections on the default port. This value must be between 1 and 65536.")

	flag.IntVar(&loglevel, "log-level", DEFAULT_LOGLEVEL, "Choose what log level you wish to use. Any value below -1 will be ignored.")
	flag.StringVar(&logfile, "log-file", "", "Choose what log file you wish to use in addition to logging output to the console. If the file can't be created or open for writing, the program will fall back to console logging only.")

	flag.StringVar(&motd, "motd", "", "Display a message of the day for the server.")
	flag.BoolVar(&motdForceDisplay, "motd-always-display", false, "Force the message of the day to be displayed upon each connection to the server, even if it hasn't changed.")

	flag.BoolVar(&sendOrigin, "send-origin", true, "Send an origin message from every message received by a client.")

	flag.BoolVar(&Launch, "launch", true, "Launch the server.")

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
		if motdForceDisplay {
			logstr += "\r\nThe server will tell each client to display this message of the day upon each connection."
		}
		Log(LOG_DEBUG, logstr)
	}

	if motd == "" && motdForceDisplay {
		Log(LOG_INFO, "The server has been told to always display a message of the day, but no message of the day has been set. The -motd-always-display parameter will be reset to false.")
		motdForceDisplay = false
	}

	if !sendOrigin {
		Log(LOG_INFO, "The server is configured to send no origin message to other clients, which may improve performance slightly, but impact the useability of your server when the origin field is required.")
	}

	if !Launch {
		Log(LOG_INFO, "The server will not be launched. Shutting down.")
		return errors.New("Server launch parameter set to false.")
	}

	if port < 1 || port > 65536 {
		Log_error("The port specified is outside the given parameter. The port parameter must be between 1 and 65536. Unable to start server.")
		return errors.New("Invalid port number.")
	}
	portstr := strconv.Itoa(port)

	ip4l := false
	ip6l := false

	if ip4 != "" && ip4 != "0" {
		Log(LOG_DEBUG, "Starting IPV4 server on address "+ip4+", using port "+portstr+".")
		S4 = NewWithTLSConfig(ip4+":"+portstr, config)
		ip4l = true
	}
	if ip6 != "" && ip6 != "0" {
		Log(LOG_DEBUG, "Starting IPV6 server on address "+ip6+", using port "+portstr+".")
		S6 = NewWithTLSConfig(ip6+":"+portstr, config)
		ip6l = true
	}
	if !ip4l && !ip6l {
		Log(LOG_DEBUG, "Starting server on all IPV4 and IPV6 addresses using port "+portstr+".")
		SAll = NewWithTLSConfig(":"+portstr, config)
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Start() int {
	var num int = 0
	var err error
	var portstr = strconv.Itoa(port)

	if S4 != nil {
		err = S4.Listen()
		if err != nil {
			Log_error("Error listening on IPV4 address.\r\n" + err.Error())
		} else {
			Log(LOG_INFO, "Listening on IPV4 address using port "+portstr+".")
			num++
		}
	}
	if S6 != nil {
		err = S6.Listen()
		if err != nil {
			Log_error("Error listening on IPV6 address.\r\n" + err.Error())
		} else {
			Log(LOG_INFO, "Listening on IPV6 address using port "+portstr+".")
			num++
		}
	}
	if SAll != nil {
		err = SAll.Listen()
		if err != nil {
			Log_error("Error listening on all addresses.\r\n" + err.Error())
			return num
		}
		Log(LOG_INFO, "Listening on all IPV4 and IPV6 addresses using port "+portstr+".")
		num++
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
