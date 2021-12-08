package server

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	connTypeSlave  string = "slave"
	connTypeMaster string = "master"
)

type ClientChannel struct {
	sync.Mutex
	name          string
	password      string
	locked        bool
	ClientsAll    map[int]*Client
	ClientsMaster map[int]*Client
	ClientsSlave  map[int]*Client
}

func (c *ClientChannel) Lmotd(ctype, name, password string) string {
	msg := "This is a locked channel. Name: " + name + "\n"
	switch ctype {
	case connTypeSlave:
		msg += "No one will be able to control your computer"
		if c.password != "" {
			msg += " unless they authenticate with the password " + c.password
		} else {
			msg += "."
		}
	case connTypeMaster:
		if (c.password != "" && password != c.password) || (c.password == "") {
			msg += "You won't be able to control any computers connected to this channel."
		}
		if c.password == password && c.password != "" {
			msg += "You are authorized to control any computer connected to this channel. Authorized with password " + password
		}
	}
	if !c.locked {
		return ""
	} else {
		return msg
	}
}

func (c *ClientChannel) Add(client *Client, password string) {
	defer c.Unlock()
	c.Lock()
	auth := false
	client.SetAuthorized(false)
	id := client.GetID()
	connection := client.GetConnectionType()
	if c.locked {
		if password == c.password && c.password != "" {
			client.SetAuthorized(true)
			auth = true
		}
	} else {
		client.SetAuthorized(true)
		auth = true
	}
	clients := c.ClientsAll
	lmotd := c.Lmotd(connection, c.name, password)
	switch connection {
	case connTypeMaster:
		_, exists := c.ClientsMaster[id]
		if exists {
			break
		}
		c.ClientsMaster[id] = client
	case connTypeSlave:
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
	if motd != "" || lmotd != "" {
		mdb := Data{
			Type:              "motd",
			Motd:              motd,
			MotdAlwaysDisplay: motdAlwaysDisplay,
		}
		if lmotd != "" {
			if mdb.Motd == "" {
				mdb.Motd = lmotd
			} else {
				mdb.Motd = lmotd + "\n" + mdb.Motd
			}
			mdb.MotdAlwaysDisplay = true
		}
		enc, encerr = Encode(mdb)
		if encerr == nil {
			client.Send(enc)
		}
	}
	logstr := "Client " + strconv.Itoa(id) + " has joined channel " + c.name
	if connection != "" {
		logstr += " as a " + connection + ". "
		if auth {
			logstr += "This client is authorized to control other computers"
		} else {
			logstr += "This client is not authorized to control other computers"
		}
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
	case connTypeMaster:
		_, exists := c.ClientsMaster[id]
		if !exists {
			break
		}
		delete(c.ClientsMaster, id)
	case connTypeSlave:
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
	auth := client.GetAuthorized()
	c.Lock()
	switch connection {
	case connTypeMaster:
		clients = c.ClientsSlave
	case connTypeSlave:
		clients = c.ClientsMaster
	default:
		clients = c.ClientsAll
	}
	c.Unlock()
	if len(clients) == 0 {
		if connection == connTypeMaster {
			client.Send([]byte("{\"type\":\"nvda_not_connected\"}"))
		}
		return
	}
	if connection == connTypeMaster && !auth {
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

func NewClientChannel(name, password string, locked bool, client *Client) *ClientChannel {
	c := &ClientChannel{
		name:          name,
		locked:        locked,
		password:      password,
		ClientsAll:    make(map[int]*Client),
		ClientsMaster: make(map[int]*Client),
		ClientsSlave:  make(map[int]*Client),
	}
	c.Add(client, password)
	return c
}

func getChannelParams(name string) (string, string, bool) {
	password := ""
	locked := false
	fl := "lock_"
	fp := "__password__"
	li := strings.Index(name, fl)
	pi := strings.Index(name, fp)
	if li == -1 && pi == -1 {
		return name, password, locked
	}
	if li == 0 {
		name = name[len(fl):]
		locked = true
		pi = strings.Index(name, fp)
	}
	if pi > 0 {
		password = name[(pi + len(fp)):]
		name = name[:pi]
		locked = true
	}
	return name, password, locked
}
