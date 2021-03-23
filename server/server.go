package server

import (
	"strconv"
	"sync"
)

var sl sync.Mutex
var EndMessage byte = '\n'
var lastID = 0
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
	Log("Client "+strconv.Itoa(client.GetID())+" has connected from "+c.GetIP(), LOG_CONNECTION)
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
		Log("This client is already disconnected. No client object was found for the closing connection, number "+strconv.Itoa(c.GetID())+".", LOG_DEBUG)
		return
	}
	sl.Lock()
	defer sl.Unlock()
	Log("Client "+strconv.Itoa(client.GetID())+" has disconnected.", LOG_CONNECTION)
	delete(clients, c)
	if len(clients) == 0 {
		clients = nil
		Log("There are now no clients connected to the server.", LOG_DEBUG)
	}
}

func AddChannel(name string, c *Client) {
	sl.Lock()
	defer sl.Unlock()
	if channels == nil {
		channels = make(map[string]*ClientChannel)
	}
	Log("Channel "+name+" has been created.", LOG_CHANNEL)
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
	Log("Channel "+name+" has been removed.", LOG_CHANNEL)
	if len(channels) == 0 {
		channels = nil
		Log("There are now no channels on the server.", LOG_DEBUG)
	}
}

var ClientConnected = func(c *Connection) {
	AddClient(c)
}

var ClientDisconnected = func(c *Connection) {
	client := FindClient(c)
	if client != nil {
		cc := client.GetChannel()
		if cc != nil {
			cc.Remove(client)
		}
	}
	RemoveClient(c)
}

var MessageReceived = func(c *Connection, pmsg []byte) {
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
		pmsg, err = JsonAdd(pmsg, "origin", id)
		if err != nil {
			cc.SendOthers(pmsg, client)
			return
		}
		cc.SendAll(pmsg, client)
		return
	}
	authErr := Authorize(client, pmsg)
	if authErr != nil {
		Log("Authorization failure for client "+strconv.Itoa(id)+".\r\n"+authErr.Error(), LOG_DEBUG)
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
