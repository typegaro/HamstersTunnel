package command

type Command struct {
	Command    string            `json:"command"`
	Parameters map[string]string `json:"parameters"`
}

type ServiceCommand struct {
	Command string `json:"command"`
	Id      string `json:"id "`
	Remote  bool   `json:"remote"`
}

type ListCommand struct {
	Command  string `json:"command"`
	Inactive bool   `json:"inactive"`
}

// NewServiceCommand rappresenta il comando 'new' per creare un nuovo servizio
type NewServiceCommand struct {
	Command       string   `json:"command"`
	ServiceName   string   `json:"service_name"`
	TCP           string   `json:"tcp_port"`  //local service port
	UDP           string   `json:"udp_port"`  //local service port
	HTTP          string   `json:"http_port"` //local service port
	RemoteIP      string   `json:"remote_ip"`
	Save          bool     `json:"save"`
	PortBlackList []string `json:"port_black_list"`
	PortWitheList []string `json:"port_withe_list"`
}
