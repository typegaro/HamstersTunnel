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

type CachedService struct {
    Ip   string `json:"ip"`
    Name string `json:"name"`
    Id   string `json:"id"`
    HTTP string `json:"http_proxy_port"`
    TCP  string `json:"tcp_proxy_port"`
    UDP  string `json:"udp_proxy_port"`
    Active bool  `json:"active"`
}
