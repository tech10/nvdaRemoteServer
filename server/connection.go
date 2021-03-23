/*
MIT License

Copyright (c) 2016 firstrow@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Connection holds info about server connection
type Connection struct {
	sync.Mutex
	conn              net.Conn
	messageTerminator byte
	ip                string
	id                int
	Server            *Server
	ctx               context.Context
	Close             context.CancelFunc
	t                 *time.Ticker
}

// TCP server
type Server struct {
	sync.Mutex
	sync.WaitGroup
	address           string // Address to open connection: localhost:9999
	config            *tls.Config
	onNewConnection   func(c *Connection)
	onClientClosed    func(c *Connection)
	onNewMessage      func(c *Connection, message []byte)
	messageTerminator byte
	ctx               context.Context
	Stop              context.CancelFunc
}

var mctx context.Context
var StopServers context.CancelFunc
var msl sync.Mutex
var ping_msg = []byte(`{"type":"ping"}`)

func init() {
	mctx, StopServers = context.WithCancel(context.Background())
}

// Read client data from channel
func (c *Connection) listen() {
	c.Lock()
	c.t = time.NewTicker(120 * time.Second)
	reader := bufio.NewReader(c.conn)
	c.Unlock()
	// Stopping and pinging our client
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				Log("Client "+strconv.Itoa(c.GetID())+" has received a signal to close.", LOG_DEBUG)
				msl.Lock()
				c.Server.Lock()
				c.t.Stop()
				c.conn.Close()
				c.Server.onClientClosed(c)
				c.Server.Unlock()
				msl.Unlock()
				c.Server.Done()
				return
			case <-c.t.C:
				err := c.Send(ping_msg)
				if err != nil {
					go c.Close()
				}
			}
		}
	}()
	c.Server.Lock()
	EndMessage := c.Server.messageTerminator
	c.Server.Unlock()
	for {
		message, err := reader.ReadBytes(EndMessage)
		if err != nil {
			Log("Error receiving message from client "+strconv.Itoa(c.GetID())+".\r\n"+err.Error()+"\r\nClosing connection.", LOG_DEBUG)
			c.Close()
			return
		}
		message = bytes.TrimSuffix(message, []byte{EndMessage})
		if len(message) == 0 {
			Log("Received empty message from client "+strconv.Itoa(c.GetID()), LOG_DEBUG)
			continue
		}
		c.t.Reset(120 * time.Second)
		Log("Data received from client "+strconv.Itoa(c.GetID())+"\r\n"+string(message), LOG_PROTOCOL)
		c.Server.onNewMessage(c, message)
	}
}

// Send bytes to client
func (c *Connection) Send(b []byte) error {
	c.Lock()
	EndMessage := c.messageTerminator
	c.Unlock()
	if len(b) == 0 {
		return nil
	}
	Log("Data sending to client "+strconv.Itoa(c.GetID())+"\r\n"+string(b), LOG_PROTOCOL)
	b = append(b, EndMessage)
	num, err := c.conn.Write(b)
	if err != nil {
		Log("Error sending message to client "+strconv.Itoa(c.GetID())+".\r\n"+err.Error()+"\r\nClosing connection.", LOG_DEBUG)
		c.Close()
		return err
	}
	if num < len(b) {
		Log("Error sending data to client "+strconv.Itoa(c.GetID())+". There were "+strconv.Itoa(num)+" bytes sent to the client, but the client should have been sent "+strconv.Itoa(len(b))+" bytes sent. Closing connection.", LOG_DEBUG)
		c.Close()
		return errors.New("Too few bytes sent to client.")
	}
	return nil
}

// Get client IP
func (c *Connection) GetIP() string {
	c.Lock()
	defer c.Unlock()
	return c.ip
}

// Get client ID
func (c *Connection) GetID() int {
	c.Lock()
	defer c.Unlock()
	return c.id
}

// Set client ID
func (c *Connection) SetID(id int) {
	c.Lock()
	defer c.Unlock()
	c.id = id
}

// Called right after server starts listening new client
func (s *Server) OnNewConnection(callback func(c *Connection)) {
	s.Lock()
	defer s.Unlock()
	s.onNewConnection = callback
}

// Called right after connection closed
func (s *Server) OnClientClosed(callback func(c *Connection)) {
	s.Lock()
	defer s.Unlock()
	s.onClientClosed = callback
}

// Called when Connection receives new message
func (s *Server) OnNewMessage(callback func(c *Connection, message []byte)) {
	s.Lock()
	defer s.Unlock()
	s.onNewMessage = callback
}

// Set message terminator
func (s *Server) MessageTerminator(terminator byte) {
	s.Lock()
	defer s.Unlock()
	s.messageTerminator = terminator
}

// Listen starts network server
func (s *Server) Listen() error {
	s.Lock()
	var listener net.Listener
	var err error
	config := s.config
	address := s.address
	s.Unlock()
	if config == nil {
		listener, err = net.Listen("tcp", address)
	} else {
		listener, err = tls.Listen("tcp", address, config)
	}
	if err != nil {
		return err
	}
	s.Lock()
	s.ctx, s.Stop = context.WithCancel(mctx)
	s.Add(1)
	s.Unlock()
	go s.accept(listener)
	return err
}

func (s *Server) accept(listener net.Listener) {
	// Stopping our server.
	go func() {
		<-s.ctx.Done()
		listener.Close()
		s.Done()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.Stop()
			break
		}
		msl.Lock()
		s.Lock()
		client := &Connection{
			conn:              conn,
			ip:                getIP(conn),
			Server:            s,
			messageTerminator: s.messageTerminator,
		}
		client.ctx, client.Close = context.WithCancel(s.ctx)
		s.Add(1)
		s.onNewConnection(client)
		s.Unlock()
		msl.Unlock()
		go client.listen()
	}
}

// Creates new tcp server instance
func New(address string) *Server {
	server := &Server{
		address:           address,
		messageTerminator: '\n',
	}

	server.OnNewConnection(ClientConnected)
	server.OnNewMessage(MessageReceived)
	server.OnClientClosed(ClientDisconnected)

	return server
}

func NewWithTLS(address, certFile, keyFile string) (*Server, error) {
	var config *tls.Config
	if certFile != "" && keyFile != "" {
		cert, cerr := tls.LoadX509KeyPair(certFile, keyFile)
		if cerr != nil {
			return nil, cerr
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	} else {
		var gerr error
		config, gerr = gen_cert()
		if gerr != nil {
			return nil, gerr
		}
	}
	server := New(address)
	server.config = config
	return server, nil
}

// Generate a self-signed certificate as long as the server is running
func serial_number() *big.Int {
	serial_num, serial_err := rand.Int(rand.Reader, big.NewInt(999999999999))
	if serial_err != nil {
		return big.NewInt(345098734305)
	}
	return serial_num
}

func gen_cert() (*tls.Config, error) {
	var ca = &x509.Certificate{
		SerialNumber: serial_number(),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"NVDARemote Server"},
			CommonName:   "Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, cerr := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)
	if cerr != nil {
		return nil, cerr
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}
	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	if err != nil {
		return nil, err
	}

	serverCert, serr := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if serr != nil {
		return nil, serr
	}

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	//serverTLSConf.InsecureSkipVerify = true

	return serverTLSConf, nil
}

func NewWithTLSConfig(address string, config *tls.Config) *Server {
	server := New(address)
	server.config = config
	return server
}

func getIP(c net.Conn) string {
	ip, _, err := net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		return ""
	}
	return strings.Trim(ip, "[]")
}
