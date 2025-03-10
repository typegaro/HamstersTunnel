package command

// Command rappresenta una struttura di comando generica
type Command struct {
	Command    string            `json:"command"`
	Parameters map[string]string `json:"parameters"`
}

// StatusCommand rappresenta il comando 'status'
type StatusCommand struct {
	Command string `json:"command"`
}

// SetStatusCommand rappresenta il comando 'set-status'
type SetStatusCommand struct {
	Command string `json:"command"`
	Value   string `json:"value"`
}

// NewServiceCommand rappresenta il comando 'new' per creare un nuovo servizio
type NewServiceCommand struct {
	Command       string `json:"command"`
	ServiceName   string `json:"service_name"`
	LocalPort     string `json:"local_port"`
	RemoteIP      string `json:"remote_ip"`
    Save          string `json:"save"`
    PortBlackList []string `json:"port_black_list"`
    PortWitheList []string `json:"port_withe_list"`
}

