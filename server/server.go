package server

import (
	"strconv"
	"sync"
)

var sl sync.Mutex
var EndMessage byte = '\n'
var lastID int = 0
var clients map[*Connection]*Client
var channels map[string]*ClientChannel

func AddClient(c *Connection) {
	sl.Lock()
	defer sl.Unlock()
	lastID++
	client := &Client{
		conn:           c,
		connectionType: "",
		id:             lastID,
		c:              nil,
	}
	c.SetID(lastID)
	if clients == nil {
		clients = make(map[*Connection]*Client)
	}
	clients[c] = client
	Log(LOG_CONNECTION, "Client "+strconv.Itoa(client.GetID())+" has connected from "+c.GetIP())
}

func FindClient(c *Connection) *Client {
	sl.Lock()
	defer sl.Unlock()
	if clients == nil {
		return nil
	}
	client, exists := clients[c]
	if !exists {
		return nil
	}
	return client
}

func RemoveClient(c *Connection) {
	client := FindClient(c)
	if client == nil {
		Log(LOG_DEBUG, "This client is already disconnected. No client object was found for the closing connection, number "+strconv.Itoa(c.GetID())+".")
		return
	}
	sl.Lock()
	defer sl.Unlock()
	Log(LOG_CONNECTION, "Client "+strconv.Itoa(client.GetID())+" has disconnected.")
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

func ClientConnected(c *Connection) {
	AddClient(c)
}

func ClientDisconnected(c *Connection) {
	client := FindClient(c)
	if client != nil {
		cc := client.GetChannel()
		if cc != nil {
			cc.Remove(client)
		}
	}
	RemoveClient(c)
}

func MessageReceived(c *Connection, pmsg []byte) {
	var err error
	client := FindClient(c)
	if client == nil {
		Log_error("A client object was not found from the connection receiving a message, number " + strconv.Itoa(c.GetID()) + ". Unexpected behavior encountered. Closing connection.")
		c.Close()
		return
	}
	id := client.GetID()
	cc := client.GetChannel()
	if cc != nil {
		if sendOrigin {
			pmsg, err = JsonAdd(pmsg, "origin", id)
			if err != nil {
				Log(LOG_DEBUG, "Error adding origin to message from client "+strconv.Itoa(id)+".\r\n"+err.Error()+"\r\nSending to all clients without origin field.")
			}
		}
		cc.SendOthers(pmsg, client)
		return
	}
	authErr := Authorize(client, pmsg)
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
