package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

var ping_msg = []byte(`{"type":"ping"}`)

const ping_sec int = 120

const write_sec int = 8

type Client struct {
	sync.Mutex
	conn              net.Conn
	messageTerminator byte
	connectionType    string
	id                int
	version           int
	ip                string
	c                 *ClientChannel
	ctx               context.Context
	Close             context.CancelFunc
	t                 *time.Ticker
	s                 *Server
	closed            bool
}

func (c *Client) ClearChannel() {
	defer c.Unlock()
	c.Lock()
	c.c = nil
}

func (c *Client) SetChannel(clientChannel *ClientChannel) {
	defer c.Unlock()
	c.Lock()
	c.c = clientChannel
}

func (c *Client) GetChannel() *ClientChannel {
	defer c.Unlock()
	c.Lock()
	return c.c
}

func (c *Client) GetID() int {
	defer c.Unlock()
	c.Lock()
	return c.id
}

func (c *Client) SetID(id int) {
	defer c.Unlock()
	c.Lock()
	c.id = id
}

func (c *Client) GetIP() string {
	defer c.Unlock()
	c.Lock()
	return c.ip
}

func (c *Client) GetConnectionType() string {
	defer c.Unlock()
	c.Lock()
	return c.connectionType
}

func (c *Client) SetConnectionType(ctype string) {
	defer c.Unlock()
	c.Lock()
	c.connectionType = ctype
}

func (c *Client) GetVersion() int {
	defer c.Unlock()
	c.Lock()
	return c.version
}

func (c *Client) SetVersion(version int) {
	defer c.Unlock()
	c.Lock()
	c.version = version
}

// Read client data from channel
func (c *Client) listen() {
	c.Lock()
	c.t = time.NewTicker(time.Duration(ping_sec) * time.Second)
	reader := bufio.NewReader(c.conn)
	EndMessage := c.messageTerminator
	c.Unlock()
	// Stopping and pinging our client
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				msl.Lock()
				c.s.Lock()
				c.t.Stop()
				c.Lock()
				c.conn.Close()
				c.closed = true
				c.Unlock()
				c.s.Unlock()
				msl.Unlock()
				return
			case <-c.t.C:
				c.Send(ping_msg)
			}
		}
	}()
	defer c.s.Done()
	defer RemoveClient(c)
	defer c.Close()
	idstr := strconv.Itoa(c.GetID())
	for {
		message, err := reader.ReadBytes(EndMessage)
		if err != nil {
			msl.Lock()
			if !stoppingServers {
				c.Lock()
				if !c.closed {
					if !errors.Is(err, io.EOF) {
						Log(LOG_DEBUG, "Error receiving message from client "+idstr+".\r\n"+err.Error()+"\r\nClosing connection.")
					}
				}
				c.Unlock()
			}
			msl.Unlock()
			return
		}
		if len(message) == 1 {
			Log(LOG_DEBUG, "Received empty message from client "+idstr)
			continue
		}
		c.t.Reset(time.Duration(ping_sec) * time.Second)
		message = bytes.TrimRight(message, string(EndMessage))
		Log(LOG_PROTOCOL, "Data received from client "+idstr+"\r\n"+string(message))
		MessageReceived(c, message)
	}
}

// Send bytes to client
func (c *Client) Send(b []byte) {
	c.Lock()
	EndMessage := c.messageTerminator
	if c.closed {
		c.Unlock()
		return
	}
	c.Unlock()
	if len(b) == 0 {
		return
	}
	Log(LOG_PROTOCOL, "Data sent to client "+strconv.Itoa(c.GetID())+"\r\n"+string(b))
	_ = c.conn.SetWriteDeadline(time.Now().Add(time.Duration(write_sec) * time.Second))
	num, err := c.conn.Write(append(b, EndMessage))
	if err != nil {
		Log(LOG_DEBUG, "Error sending message to client "+strconv.Itoa(c.GetID())+".\r\n"+err.Error()+"\r\nClosing connection.")
		c.Close()
		return
	}
	if num < len(b)+1 {
		Log(LOG_DEBUG, "Error sending data to client "+strconv.Itoa(c.GetID())+". There were "+strconv.Itoa(num)+" bytes sent to the client, but the client should have been sent "+strconv.Itoa(len(b)+1)+" bytes sent. Closing connection.")
		c.Close()
		return
	}
}
