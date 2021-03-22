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

var loglevel int

const DEFAULT_LOGLEVEL int = 0
const LOG_INFO int = 0
const LOG_CONNECTION int = 1
const LOG_CHANNEL int = 2
const LOG_DEBUG int = 3
const LOG_PROTOCOL int = 4

var log_standard *log.Logger
var log_error *log.Logger

var S4 *Server
var S6 *Server
var SAll *Server

func Configure() error {
	flag.StringVar(&Cert, "cert", "", "SSL certificate file to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")
	flag.StringVar(&Key, "key", "", "SSL key to use for the server's TLS connection, must point to an existing file. If this is empty, the server will automatically generate its own self-signed certificate.")

	flag.StringVar(&ip4, "ip4", "", "IPV4 address for the server to listen for connections on. This can be blank if desired, in which case, the server will listen on all IPV4 addresses.")
	flag.StringVar(&ip6, "ip6", "", "IPV6 address for the server to listen for connections on. This can be blank if desired, in which case, the server will listen on all IPV6 addresses.")
	flag.IntVar(&port, "port", DEFAULT_PORT, "The port that the server will listen for connections on. This can be blank if desired, in which case, the server will listen for connections on the default port, "+strconv.Itoa(DEFAULT_PORT)+". This value must be between 1 and 65536.")

	flag.IntVar(&loglevel, "loglevel", DEFAULT_LOGLEVEL, "Choose what log level you wish to use. Any value below -1 will be ignored.")
	flag.Parse()

	log_standard = log.New(os.Stdout, "", log.LstdFlags)
	log_error = log.New(os.Stderr, "[ERROR]:", log.LstdFlags)

	Log("Initializing configuration.", LOG_INFO)

	generate := false
	var config *tls.Config
	var err error

	if Cert != "" && !fileExists(Cert) {
		Log("The certificate file at "+Cert+" does not exist.", LOG_INFO)
		generate = true
	}
	if Key != "" && !fileExists(Key) {
		Log("The key file at "+Key+" does not exist.", LOG_INFO)
		generate = true
	}
	if Cert == "" || Key == "" {
		generate = true
	}
	if port < 1 || port > 65536 {
		Log_error("The port specified is outside the given parameter. The port parameter must be between 1 and 65536. Unable to start server.")
		return errors.New("Invalid port number.")
	}

	if generate {
		Log("Attempting to generate self-signed SSL certificate.", LOG_DEBUG)
		config, err = gen_cert()
		if err != nil {
			Log_error("Unable to generate self-signed certificate.\r\n" + err.Error() + "\r\nUnable to start server.")
			return err
		}
		Log("SSL certificate generated.", LOG_DEBUG)
	} else {
		cert, cerr := tls.LoadX509KeyPair(Cert, Key)
		if cerr != nil {
			Log_error("Error loading certificate and key files.\r\n" + cerr.Error() + "\r\nUnable to start server.")
			return cerr
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	portstr := strconv.Itoa(port)

	ip4l := false
	ip6l := false

	if ip4 != "" && ip4 != "0" {
		Log("Starting IPV4 server on address "+ip4+", listening on port "+portstr, LOG_DEBUG)
		S4 = NewWithTLSConfig(ip4+":"+portstr, config)
	}
	if ip6 != "" && ip6 != "0" {
		Log("Starting IPV6 server on address "+ip6+", listening on port "+portstr, LOG_DEBUG)
		S6 = NewWithTLSConfig(ip6+":"+portstr, config)
	}
	if !ip4l && !ip6l {
		Log("Listening on all IPV4 and IPV6 addresses using port "+portstr, LOG_DEBUG)
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

	if S4 != nil {
		err = S4.Listen()
		if err != nil {
			Log_error("Error listening on IPV4 address.\r\n" + err.Error())
		} else {
			Log("Listening on IPV4 address.", LOG_INFO)
			num++
		}
	}
	if S6 != nil {
		err = S6.Listen()
		if err != nil {
			Log_error("Error listening on IPV6 address.\r\n" + err.Error())
		} else {
			Log("Listening on IPV6 address.", LOG_INFO)
			num++
		}
	}
	if SAll != nil {
		err = SAll.Listen()
		if err != nil {
			Log_error("Error listening on all addresses.\r\n" + err.Error())
			return num
		}
		Log("Listening on all IPV4 and IPV6 addresses.", LOG_INFO)
		num++
	}
	Log("Number of servers started: "+strconv.Itoa(num), LOG_DEBUG)
	go signals_init()
	return num
}
