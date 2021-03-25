package server

import (
	"sort"
	"strconv"
	"sync"
)

type ClientChannel struct {
	sync.Mutex
	name          string
	ClientsAll    map[int]*Client
	ClientsMaster map[int]*Client
	ClientsSlave  map[int]*Client
}

func (c *ClientChannel) Add(client *Client) {
	defer c.Unlock()
	c.Lock()
	id := client.GetID()
	connection := client.GetConnectionType()
	clients := c.ClientsAll
	switch connection {
	case "master":
		_, exists := c.ClientsMaster[id]
		if exists {
			break
		}
		c.ClientsMaster[id] = client
	case "slave":
		_, exists := c.ClientsSlave[id]
		if exists {
			break
		}
		c.ClientsSlave[id] = client
	}
	_, exists := c.ClientsAll[id]
	if exists {
		return
	}
	c.ClientsAll[id] = client
	client.SetChannel(c)
	scdb := Data{
		Type:    "client_joined",
		Channel: c.name,
		ID:      id,
		Client: &ClientData{
			ID:             id,
			ConnectionType: connection,
		},
	}
	enc, encerr := Encode(scdb)
	if encerr == nil {
		c.SendAll(enc, client)
	} else {
		Log(LOG_DEBUG, "Error encoding JSON for client "+strconv.Itoa(id)+" while trying to add them to channel "+c.name+"\r\n"+encerr.Error())
	}

	scdb.Type = "channel_joined"
	scdb.Origin = id
	scdb.ID = 0
	scdb.Client = nil
	scdb.Motd = motd
	scdb.MotdForceDisplay = motdForceDisplay
	if len(clients) > 0 {
		scdb.UserIds = make([]int, 0, len(clients))
		scdb.Clients = make([]ClientData, 0, len(clients))
		var ctype string
		for cid, cc := range clients {
			if cid == id {
				continue
			}
			ctype = cc.GetConnectionType()
			scdb.UserIds = append(scdb.UserIds, cid)
			scdb.Clients = append(scdb.Clients, ClientData{
				ID:             cid,
				ConnectionType: ctype,
			})
		}
		if len(scdb.UserIds) == 0 {
			scdb.UserIds = nil
			scdb.Clients = nil
		} else if len(scdb.UserIds) > 1 {
			sort.Ints(scdb.UserIds)
			sort.SliceStable(scdb.Clients,
				func(i, j int) bool {
					return scdb.Clients[i].ID < scdb.Clients[j].ID
				})
		}

	}
	enc, encerr = Encode(scdb)
	if encerr == nil {
		client.Send(enc)
	}
	logstr := "Client " + strconv.Itoa(id) + " has joined channel " + c.name
	if connection != "" {
		logstr += " as a " + connection
	}
	Log(LOG_CHANNEL, logstr+".")
}

func (c *ClientChannel) Remove(client *Client) {
	defer c.EndIfEmpty()
	defer c.Unlock()
	c.Lock()
	id := client.GetID()
	connection := client.GetConnectionType()
	switch connection {
	case "master":
		_, exists := c.ClientsMaster[id]
		if !exists {
			break
		}
		delete(c.ClientsMaster, id)
	case "slave":
		_, exists := c.ClientsSlave[id]
		if !exists {
			break
		}
		delete(c.ClientsSlave, id)
	}
	_, exists := c.ClientsAll[id]
	if exists {
		delete(c.ClientsAll, id)
	}
	client.ClearChannel()
	scdb := Data{
		Type:   "client_left",
		ID:     id,
		Origin: id,
		Client: &ClientData{
			ID:             id,
			ConnectionType: connection,
		},
	}
	enc, encerr := Encode(scdb)
	if encerr == nil {
		c.SendAll(enc, client)
	}
	Log(LOG_CHANNEL, "Client "+strconv.Itoa(id)+" has left channel "+c.name)
}

func (c *ClientChannel) EndIfEmpty() bool {
	c.Lock()
	if len(c.ClientsAll) > 0 {
		c.Unlock()
		return false
	}
	c.Unlock()
	c.Quit()
	return true
}

func (c *ClientChannel) Quit() {
	defer c.Unlock()
	c.Lock()
	if c.ClientsAll == nil {
		RemoveChannel(c.name)
		return
	}
	for id, client := range c.ClientsAll {
		delete(c.ClientsMaster, id)
		delete(c.ClientsSlave, id)
		client.ClearChannel()
	}
	c.ClientsAll = nil
	c.ClientsMaster = nil
	c.ClientsSlave = nil
	RemoveChannel(c.name)
}

func (c *ClientChannel) SendAll(msg []byte, client *Client) {
	if c.ClientsAll == nil || len(c.ClientsAll) == 0 {
		return
	}
	for _, sc := range c.ClientsAll {
		if client != nil && client == sc {
			continue
		}
		sc.Send(msg)
	}
}

func (c *ClientChannel) SendOthers(msg []byte, client *Client) {
	if client == nil {
		return
	}
	connection := client.GetConnectionType()
	var clients map[int]*Client
	c.Lock()
	switch connection {
	case "master":
		clients = c.ClientsSlave
	case "slave":
		clients = c.ClientsMaster
	default:
		clients = c.ClientsAll
	}
	c.Unlock()
	if len(clients) == 0 {
		if connection == "master" {
			client.Send([]byte("{\"type\":\"nvda_not_connected\"}"))
		}
		return
	}
	for _, sc := range clients {
		if sc == client {
			continue
		}
		sc.Send(msg)
	}
}

func (c *ClientChannel) Name() string {
	c.Lock()
	defer c.Unlock()
	return c.name
}

func NewClientChannel(name string, client *Client) *ClientChannel {
	c := &ClientChannel{
		name:          name,
		ClientsAll:    make(map[int]*Client),
		ClientsMaster: make(map[int]*Client),
		ClientsSlave:  make(map[int]*Client),
	}
	c.Add(client)
	return c
}

type Client struct {
	sync.Mutex
	conn           *Connection
	connectionType string
	id             int
	version        int
	c              *ClientChannel
}

func (c *Client) Close() {
	defer c.Unlock()
	c.Lock()
	c.conn.Close()
}

func (c *Client) CloseForce() {
	c.Lock()
	defer c.Unlock()
	c.conn.conn.Close()
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

func (c *Client) Send(msg []byte) {
	if len(msg) == 0 {
		return
	}
	_ = c.conn.Send(msg)
}
