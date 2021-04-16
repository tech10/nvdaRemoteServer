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
			Log_error(v...)
		} else {
			Log(c.ll[i], v...)
		}
	}
	c.ls = nil
	c.ll = nil
	c.le = nil
}

func (c *Cfg) Write(file string) error {
	if file == "" {
		err := errors.New("An empty file is an invalid parameter. Not writing.")
		c.Log_error(err)
		return err
	}
	if c.IsDefault() {
		err := errors.New("Default parameters have been used. Nothing to write to configuration file.")
		c.Log_error(err)
		return err
	}
	d, err := cfg_write(c)
	if err != nil {
		c.Log_error("Unable to encode json for writing.\n" + err.Error())
		return err
	}
	file = fullPath(file)
	c.Log(LOG_DEBUG, "Writing to configuration file "+file)
	err = file_rewrite(file, d)
	if err != nil {
		c.Log_error(err)
		return err
	}
	c.Log(LOG_DEBUG, "Configuration file successfully written to "+file)
	return nil
}

func (c *Cfg) SearchFile(f string) string {
	f = fullPath(f)
	c.Log(LOG_DEBUG, "Searching for configuration file at "+f)
	if fileExists(f) {
		return f
	}
	if default_conf_file(confFile) {
		c.Log(LOG_DEBUG, "Failed to find configuration file.")
	} else {
		c.Log_error("The configuration file at " + f + " does not exist.")
	}
	return ""
}

func (c *Cfg) FindFile() string {
	if !default_conf_file(confFile) {
		return c.SearchFile(confFile)
	}
	cf := c.SearchFile(DEFAULT_CONF_NAME)
	if cf != "" {
		return cf
	}
	return c.SearchFile(DEFAULT_CONF_DIR + PS + DEFAULT_CONF_NAME)
}

func (c *Cfg) ReadFile(f string) ([]byte, error) {
	c.Log(LOG_DEBUG, "Reading configuration file at "+f)
	d, err := file_read(f)
	if err != nil {
		c.Log_error("Unable to read configuration file at " + f + "\n" + err.Error())
		return nil, err
	}
	c.Log(LOG_DEBUG, "Successfully read configuration file.")
	return d, nil
}

func (c *Cfg) Decode(d []byte) error {
	return cfg_read(d, c)
}

func (c *Cfg) Read() error {
	f := c.FindFile()
	if f == "" {
		return errors.New("No file found.")
	}
	d, err := c.ReadFile(f)
	if err != nil {
		return errors.New("Error reading " + f + "\n" + err.Error())
	}
	c.Log(LOG_DEBUG, "Decoding data from configuration file.")
	err = c.Decode(d)
	if err != nil {
		c.Log_error("Unable to decode, invalid data in " + f + "\n" + err.Error())
		return err
	}
	c.Log(LOG_DEBUG, "Data successfully decoded.")
	c.file = f
	return nil
}

func (c *Cfg) Setup() error {
	err := c.Read()
	if err != nil {
		if !default_conf_file(confFile) {
			return err
		}
		return nil
	}
	return nil
}
