package server

type Data struct {
	Type              string       `json:"type"`
	Channel           string       `json:"channel,omitempty"`
	ConnectionType    string       `json:"connection_type,omitempty"`
	Version           int          `json:"version,omitempty"`
	Origin            int          `json:"origin,omitempty"`
	Key               string       `json:"key,omitempty"`
	ID                int          `json:"user_id,omitempty"`
	UserIds           []int        `json:"user_ids,omitempty"`
	Clients           []ClientData `json:"clients,omitempty"`
	Client            *ClientData  `json:"client,omitempty"`
	Error             string       `json:"error,omitempty"`
	Motd              string       `json:"motd,omitempty"`
	MotdAlwaysDisplay bool         `json:"force_display,omitempty"`
}

type ClientData struct {
	ID             int    `json:"id"`
	ConnectionType string `json:"connection_type"`
}

type Cfg struct {
	PidFile           string      `json:"pid_file,omitempty"`
	LogFile           string      `json:"log_file,omitempty"`
	LogLevel          int         `json:"log_level,omitempty"`
	Addresses         AddressList `json:"addresses,omitempty"`
	Cert              string      `json:"cert_file,omitempty"`
	Key               string      `json:"key_file,omitempty"`
	Motd              string      `json:"motd,omitempty"`
	MotdAlwaysDisplay bool        `json:"motd_always_display,omitempty"`
	SendOrigin        bool        `json:"send_origin,omitempty"`
}
