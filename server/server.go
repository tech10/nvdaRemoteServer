package server

import (
	"runtime"
	"strconv"
	"sync"
)

var sl sync.Mutex
var EndMessage byte = '\n'
var lastID int = 0
var clients map[*Client]struct{}
var channels map[string]*ClientChannel

func AddClient(c *Client) {
	sl.Lock()
	defer sl.Unlock()
	lastID++
	c.SetID(lastID)
	if clients == nil {
		clients = make(map[*Client]struct{})
	}
	clients[c] = struct{}{}
	Log(LOG_CONNECTION, "Client "+strconv.Itoa(lastID)+" has connected from "+c.GetIP())
}

func FindClient(c *Client) bool {
	sl.Lock()
	defer sl.Unlock()
	if clients == nil {
		return false
	}
	_, exists := clients[c]
	return exists
}

func RemoveClient(c *Client) {
	if !FindClient(c) {
		Log(LOG_DEBUG, "Client "+strconv.Itoa(c.GetID())+" is already disconnected.")
		return
	}
	sl.Lock()
	defer sl.Unlock()
	Log(LOG_CONNECTION, "Client "+strconv.Itoa(c.GetID())+" has disconnected.")
	delete(clients, c)
	if len(clients) == 0 {
		clients = nil
		Log(LOG_DEBUG, "There are now no clients connected to the server.")
	}
}

func AddChannel(name string, c *Client) {
	sl.Lock()
	defer sl.Unlock()
	if channels == nil {
		channels = make(map[string]*ClientChannel)
	}
	Log(LOG_CHANNEL, "Channel "+name+" has been created.")
	cc := NewClientChannel(name, c)
	channels[name] = cc
}

func FindChannel(name string) *ClientChannel {
	sl.Lock()
	defer sl.Unlock()
	if channels == nil {
		return nil
	}
	c, exists := channels[name]
	if !exists {
		return nil
	}
	return c
}

func RemoveChannel(name string) {
	c := FindChannel(name)
	if c == nil {
		return
	}
	sl.Lock()
	defer sl.Unlock()
	delete(channels, name)
	Log(LOG_CHANNEL, "Channel "+name+" has been removed.")
	if len(channels) == 0 {
		channels = nil
		Log(LOG_DEBUG, "There are now no channels on the server.")
	}
}

func ClientConnected(c *Client) {
	AddClient(c)
}

func ClientDisconnected(c *Client) {
	cc := c.GetChannel()
	if cc != nil {
		cc.Remove(c)
	}
	RemoveClient(c)
}

func MessageReceived(c *Client, pmsg []byte) {
	var err error
	if !FindClient(c) {
		Log_error("A client object was not found from the connection receiving a message, number " + strconv.Itoa(c.GetID()) + ". Unexpected behavior encountered. Closing connection.")
		runtime.Goexit()
	}
	id := c.GetID()
	cc := c.GetChannel()
	if cc != nil {
		if sendOrigin {
			pmsg, err = JsonAdd(pmsg, "origin", id)
			if err != nil {
				Log(LOG_DEBUG, "Error adding origin to message from client "+strconv.Itoa(id)+".\r\n"+err.Error()+"\r\nSending to all clients without origin field.")
			}
		}
		cc.SendOthers(pmsg, c)
		return
	}
	authErr := Authorize(c, pmsg)
	if authErr != nil {
		Log(LOG_DEBUG, "Authorization failure for client "+strconv.Itoa(id)+".\r\n"+authErr.Error())
		c.Close()
		return
	}
}

func Authorize(c *Client, data []byte) error {
	decode, err := Decode(data)
	if err != nil {
		return err
	}
	return cmd_exec(c, &decode)
}
