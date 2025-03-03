package models 


type PublicService struct {
    Info ServiceInfo`json:"info"` 
	HTTP  PortPair`json:"http_port_pair"`
	TCP   PortPair `json:"tcp_port_pair"`
	UDP   PortPair `json:"udp_port_pair"`
	Active     bool   `json:"active"`
}
type ServiceInfo struct {
	Id string `json:"id"`
    Name string `json:"name"` 
}

type PortPair struct{
    External string `json:"external"`
    Internal string `json:"internal"`
}

type NewServiceReq struct {
    Name string `json:"name"` 
    TCP bool `json:"tcp"` 
    UDP bool `json:"udp"`
    HTTP bool `json:"http"`
    PortBlackList []string `json:"port_black_list"`
}

type LocalService struct {
    Info ServiceInfo`json:"info"` 
	LocalHTTP  string `json:"local_http"`
	LocalTPC   string `json:"local_tcp"` 
	LocalUDP   string `json:"local_udp"` 
    Options []string  `json:"options"` //TODO: Make an Option type
}
