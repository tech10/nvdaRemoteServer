package server

type Data struct {
	Type           string       `json:"type"`
	Channel        string       `json:"channel,omitempty"`
	ConnectionType string       `json:"connection_type,omitempty"`
	Version        int          `json:"version,omitempty"`
	Origin         int          `json:"origin,omitempty"`
	Key            string       `json:"key,omitempty"`
	ID             int          `json:"user_id,omitempty"`
	UserIds        []int        `json:"user_ids,omitempty"`
	Clients        []ClientData `json:"clients,omitempty"`
	Client         *ClientData  `json:"client,omitempty"`
	Error          string       `json:"error,omitempty"`
	Motd           string       `json:"motd,omitempty"`
	MotdForce      bool         `json:"force_display,omitempty"`
}

type ClientData struct {
	ID             int    `json:"id"`
	ConnectionType string `json:"connection_type"`
}
