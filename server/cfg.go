package server

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
