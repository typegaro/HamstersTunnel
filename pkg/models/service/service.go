package models 

type Service struct {
    Id string `json:"id"`
    Name string `json:"name"`
	Active  bool        `json:"active"`
	Options []string    `json:"options"`
	HTTP  *PortPair `json:"http_port_pair"`
	TCP   *PortPair `json:"tcp_port_pair"`
	UDP   *PortPair `json:"udp_port_pair"`
}

//In Local case:
//- Proxy is the port on the remote
//- Client is the local service port
//In Remote case:
//- Proxy is the port listening for the final user 
//- Client is the port for the client 
//-  Remote.Client == Local.Proxy
type PortPair struct {
    Proxy string `json:"proxy"`
    Client string `json:"client"` 
}

type NewServiceReq struct {
    Name string `json:"name"` 
    TCP bool `json:"tcp"` 
    UDP bool `json:"udp"`
    HTTP bool `json:"http"`
    PortBlackList []string `json:"port_black_list"`
    PortWitheList []string `json:"port_white_list"`
    Options []string  `json:"options"` 
}

type ServiceRes struct{
    Id string `json:"id"`
    HTTP  string `json:"http_port"`
	TCP   string `json:"tcp_port"`
	UDP   string `json:"udp_port"`
}
