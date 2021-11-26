/*
MIT License

Copyright (c) 2016 firstrow@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"sync"
)

// TCP server
type Server struct {
	sync.Mutex
	sync.WaitGroup
	address           string // Address to open connection: localhost:9999
	config            *tls.Config
	messageTerminator byte
	ctx               context.Context
	Stop              context.CancelFunc
}

var (
	mctx            context.Context
	StopServers     context.CancelFunc
	msl             sync.Mutex
	stoppingServers bool = false
)

func init() {
	mctx, StopServers = context.WithCancel(context.Background())
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
	s.Lock()
	address := s.address
	s.Unlock()
	// Stopping our server.
	go func() {
		<-s.ctx.Done()
		msl.Lock()
		if !stoppingServers {
			Log(LOG_DEBUG, "The server at "+address+" has received a signal to stop.")
		}
		msl.Unlock()
		listener.Close()
		s.Done()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			msl.Lock()
			if !stoppingServers {
				Log(LOG_DEBUG, "Error accepting connections on the server at "+address+"\r\n"+err.Error()+"\r\nStopping server.")
			}
			msl.Unlock()
			s.Stop()
			break
		}
		msl.Lock()
		s.Lock()
		client := &Client{
			conn:              conn,
			ip:                getIP(conn),
			s:                 s,
			messageTerminator: s.messageTerminator,
			closed:            false,
		}
		client.ctx, client.Close = context.WithCancel(s.ctx)
		s.Add(1)
		AddClient(client)
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

	return server
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
