package server

import (
	"errors"
)

type Cfg struct {
	PidFile           string      `json:"pid_file"`
	LogFile           string      `json:"log_file"`
	LogLevel          int         `json:"log_level"`
	Addresses         AddressList `json:"addresses"`
	Cert              string      `json:"cert_file"`
	Key               string      `json:"key_file"`
	Motd              string      `json:"motd"`
	MotdAlwaysDisplay bool        `json:"motd_always_display"`
	SendOrigin        bool        `json:"send_origin"`
	ll                []int
	ls                [][]interface{}
	le                []bool
	file              string
}

func cfg_default() *Cfg {
	return &Cfg{
		PidFile:           DEFAULT_PID_FILE,
		LogFile:           DEFAULT_LOG_FILE,
		LogLevel:          DEFAULT_LOG_LEVEL,
		Addresses:         AddressList{DEFAULT_ADDRESS},
		Cert:              DEFAULT_CERT_FILE,
		Key:               DEFAULT_KEY_FILE,
		Motd:              DEFAULT_MOTD,
		MotdAlwaysDisplay: DEFAULT_MOTD_ALWAYS_DISPLAY,
		SendOrigin:        DEFAULT_SEND_ORIGIN,
		ll:                make([]int, 0),
		ls:                make([][]interface{}, 0),
		le:                make([]bool, 0),
	}
}

func (c *Cfg) IsDefault() bool {
	if !default_pid_file(c.PidFile) {
		return false
	}
	if !default_log_file(c.LogFile) {
		return false
	}
	if !default_log_level(c.LogLevel) {
		return false
	}
	if !default_addresses(c.Addresses) {
		return false
	}
	if !default_cert_file(c.Cert) {
		return false
	}
	if !default_key_file(c.Key) {
		return false
	}
	if !default_motd(c.Motd) {
		return false
	}
	if !default_motd_always_display(c.MotdAlwaysDisplay) {
		return false
	}
	if !default_send_origin(c.SendOrigin) {
		return false
	}
	return true
}

func (c *Cfg) Log(level int, msg ...interface{}) {
	if default_conf_file(confFile) {
		level = LOG_DEBUG
	}
	c.ls = append(c.ls, msg)
	c.ll = append(c.ll, level)
	c.le = append(c.le, false)
}

func (c *Cfg) Log_error(msg ...interface{}) {
	c.ls = append(c.ls, msg)
	c.ll = append(c.ll, LOG_SILENT)
	c.le = append(c.le, true)
}

func (c *Cfg) LogWrite() {
	if c.ls == nil {
		return
	}
	for i, v := range c.ls {
		if c.le[i] {
			Log_error(v)
		} else {
			Log(c.ll[i], v)
		}
	}
	c.ls = nil
	c.ll = nil
	c.le = nil
}

func (c *Cfg) Write(file string) error {
	if file == "" {
		err := errors.New("An empty file is an invalid parameter. Not writing.")
		c.Log(LOG_DEBUG, err)
		return err
	}
	if c.IsDefault() {
		err := errors.New("Default parameters have been used. Nothing to write to configuration file.")
		c.Log(LOG_DEBUG, err)
		return err
	}
	d, err := cfg_write(c)
	if err != nil {
		c.Log(LOG_DEBUG, "Unable to encode json for writing.\n"+err.Error())
		return err
	}
	file = fullPath(file)
	c.Log(LOG_DEBUG, "Writing to configuration file "+file)
	err = file_rewrite(c.file, d)
	if err != nil {
		c.Log_error(err)
		return err
	}
	c.Log(LOG_DEBUG, "Configuration file successfully written to "+file)
	return nil
}
