package server

import (
	"errors"
	"strconv"
)

var command = make(map[string]func(*Client, *Data))

func cmd_exists(cmd string) bool {
	_, exists := command[cmd]
	return exists
}

func cmd_add(cmd string, cfunc func(*Client, *Data)) {
	command[cmd] = cfunc
}

func cmd_exec(c *Client, db *Data) error {
	cmd := db.Type
	if !cmd_exists(cmd) {
		return errors.New("The command " + cmd + " does not exist.")
	}
	command[cmd](c, db)
	return nil
}

func init() {
	cmd_add("join", func(c *Client, db *Data) {
		if c.GetChannel() != nil {
			enc, encerr := Encode(Data{
				Type:  "error",
				Error: "already_joined",
			})
			if encerr == nil {
				c.Send(enc)
				return
			} else {
				Log("JSON encoding error for client "+strconv.Itoa(c.GetID()), LOG_DEBUG)
				return
			}
		}
		if db.Channel == "" {
			enc, encerr := Encode(Data{
				Type:  "error",
				Error: "invalid_parameters",
			})
			if encerr == nil {
				c.Send(enc)
				return
			} else {
				Log("JSON encoding error for client "+strconv.Itoa(c.GetID()), LOG_DEBUG)
				return
			}
		}

		c.SetConnectionType(db.ConnectionType)
		cc := FindChannel(db.Channel)
		if cc != nil {
			cc.Add(c)
			return
		}
		AddChannel(db.Channel, c)
	})

	cmd_add("protocol_version", func(c *Client, db *Data) {
		if db.Version <= 0 {
			Log("Client "+strconv.Itoa(c.GetID())+" has tried to register an invalid version number.", LOG_DEBUG)
			enc, encerr := Encode(Data{
				Type:  "error",
				Error: "invalid_parameters",
			})
			if encerr == nil {
				c.Send(enc)
				return
			} else {
				Log("JSON encoding error for client "+strconv.Itoa(c.GetID()), LOG_DEBUG)
				return
			}
		}
		c.SetVersion(db.Version)
		Log("Client "+strconv.Itoa(c.GetID())+" has set protocol version "+strconv.Itoa(db.Version)+".", LOG_DEBUG)
	})

	cmd_add("generate_key", func(c *Client, db *Data) {
		key := gen_key()
		enc, encerr := Encode(Data{
			Type: "generate_key",
			Key:  key,
		})
		if encerr != nil {
			Log("JSON encoding error for client "+strconv.Itoa(c.GetID()), LOG_DEBUG)
			return
		}
		c.Send(enc)
		Log("Client "+strconv.Itoa(c.GetID())+" has generated a key: "+key, LOG_DEBUG)
		c.Close()
	})

}
