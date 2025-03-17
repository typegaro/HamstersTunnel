package models

type ServerService struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Active  bool            `json:"active"`
	HTTP    *ServerPortPair `json:"http_port_pair"`
	TCP     *ServerPortPair `json:"tcp_port_pair"`
	UDP     *ServerPortPair `json:"udp_port_pair"`
	Options []string        `json:"options"`
}

type ServerPortPair struct {
	Proxy        string `json:"proxy"`
	Client       string `json:"client"`
	Iniitialized bool   `json:"iniitialized"`
}

type ClientService struct {
	Id      string          `json:"id"`
	Ip      string          `json:"ip"`
	Name    string          `json:"name"`
	Active  bool            `json:"active"`
	HTTP    *ClientPortPair `json:"http_port_pair"`
	TCP     *ClientPortPair `json:"tcp_port_pair"`
	UDP     *ClientPortPair `json:"udp_port_pair"`
	Options []string        `json:"options"`
}

type ClientPortPair struct {
	Remote       string `json:"remote"`
	Local        string `json:"local"`
	Iniitialized bool   `json:"iniitialized"`
}

type NewServiceReq struct {
	Name          string   `json:"name"`
	TCP           bool     `json:"tcp"`
	UDP           bool     `json:"udp"`
	HTTP          bool     `json:"http"`
	PortBlackList []string `json:"port_black_list"`
	PortWitheList []string `json:"port_white_list"`
	Options       []string `json:"options"`
}

type ServiceRes struct {
	Id   string `json:"id"`
	HTTP string `json:"http_port"`
	TCP  string `json:"tcp_port"`
	UDP  string `json:"udp_port"`
}
