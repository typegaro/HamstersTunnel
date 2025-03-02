package models 

type Service struct{
	Name       string `json:"name"`
    PService PublicService
    LService LocalService
}

type PublicService struct {
	PublicHTTP int    `json:"public_http"`
	PublicTCP  int    `json:"public_tcp"`
	PublicUDP  int    `json:"public_udp"`
}
type LocalService struct {
	LocalHTTP  string `json:"local_http"`
	LocalTPC   string `json:"local_tcp"` 
	LocalUDP   string `json:"local_udp"` 
	Active     bool   `json:"active"`
}
