package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/typegaro/HamstersTunnel/internal/daemon/client_memory"
	"github.com/typegaro/HamstersTunnel/pkg/command"
	"github.com/typegaro/HamstersTunnel/pkg/interfaces"
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
	"github.com/typegaro/HamstersTunnel/pkg/reversetunnel"
)

type Daemon struct {
	status   string
	listener net.Listener
	memory   interfaces.ClientMemory
	services map[string]*models.ClientService
}

func NewDaemon() *Daemon {
	return &Daemon{
		status:   "Running",
		memory:   &client_memory.FileSystemMemory{},
		services: make(map[string]*models.ClientService),
	}
}

func (d *Daemon) Init() {
	d.memory.Init()
	d.InitSocket()
	d.WakeUpServices()
	for {
		conn, err := d.listener.Accept()
		defer d.listener.Close()
		if err != nil {
			log.Println("Errore nella connessione:", err)
			continue
		}
		go d.handleConnection(conn)
	}
}

func getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon` // Named Pipe on Windows
	}
	return "/tmp/mydaemon.sock" // Unix socket on Linux/macOS
}

func (d *Daemon) InitSocket() {
	socketPath := getSocketPath()

	if runtime.GOOS != "windows" {
		os.Remove(socketPath)
	}

	var err error
	if runtime.GOOS == "windows" {
		d.listener, err = net.Listen("npipe", socketPath)
	} else {
		d.listener, err = net.Listen("unix", socketPath)
	}

	if err != nil {
		log.Fatal("Errore nell'apertura del socket:", err)
	}

	fmt.Println("Daemon in ascolto su", socketPath)
}

func (d *Daemon) WakeUpServices() error {
	srvs, err := d.memory.GetActiveServices()
	if err != nil {
		return err
	}

	for _, srv := range srvs {
		if srv.TCP != nil {
			go reversetunnel.StartLocalTCPTunnel(srv.TCP.Remote, srv.TCP.Local)
		}
		d.addService(srv)
	}
	return nil
}

func (d *Daemon) addService(cs *models.ClientService) error {
	if _, exists := d.services[cs.Id]; exists {
		return fmt.Errorf("service with id %s already exists", cs.Id)
	}
	d.services[cs.Id] = cs
	return nil
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Errore nella lettura:", err)
		return
	}

	var request command.Command
	err = json.Unmarshal(buffer[:n], &request)
	if err != nil {
		fmt.Println("Errore nel parsing JSON:", err)
		return
	}

	switch request.Command {
	case "status":
		var statusCommand command.StatusCommand
		if err := json.Unmarshal(buffer[:n], &statusCommand); err != nil {
			conn.Write([]byte("Errore nel parsing del comando 'status'\n"))
			return
		}
		d.handleStatus(conn)
	case "new":
		var newServiceCommand command.NewServiceCommand
		if err := json.Unmarshal(buffer[:n], &newServiceCommand); err != nil {
			conn.Write([]byte("Errore nel parsing del comando 'new'\n"))
			return
		}
		d.handleNewService(conn, newServiceCommand)
	default:
		conn.Write([]byte("Unknown command\n"))
	}
}

func (d *Daemon) handleStatus(conn net.Conn) {
	response := map[string]string{"status": d.status}
	json.NewEncoder(conn).Encode(response)
}

func (d *Daemon) handleNewService(conn net.Conn, command command.NewServiceCommand) {
	url := fmt.Sprintf("http://%s/service?save=%v", command.RemoteIP, command.Save)

	payload := models.NewServiceReq{
		Name:          command.ServiceName,
		TCP:           true,
		UDP:           true,
		HTTP:          true,
		PortBlackList: []string{},
		PortWitheList: []string{},
		Options:       []string{},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		conn.Write([]byte("Error encoding JSON: " + err.Error() + "\n"))
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		conn.Write([]byte("Error making POST request: " + err.Error() + "\n"))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		conn.Write([]byte("Error reading response body: " + err.Error() + "\n"))
		return
	}

	var serviceRes models.ServiceRes
	if err := json.Unmarshal(body, &serviceRes); err != nil {
		conn.Write([]byte("Error decoding response body: " + err.Error() + "\n"))
		return
	}

	s := &models.ClientService{
		Id:     serviceRes.Id,
		Name:   command.ServiceName,
		Ip:     command.RemoteIP,
		Active: true,
	}
	if serviceRes.TCP != "" {
		s.TCP = &models.ClientPortPair{
			Remote:       serviceRes.TCP,
			Local:        command.TCP,
			Iniitialized: true,
		}
	}
	if serviceRes.UDP != "" {
		s.UDP = &models.ClientPortPair{
			Remote:       serviceRes.UDP,
			Local:        command.UDP,
			Iniitialized: true,
		}
	}
	if serviceRes.HTTP != "" {
		s.HTTP = &models.ClientPortPair{
			Remote:       serviceRes.HTTP,
			Local:        command.HTTP,
			Iniitialized: true,
		}
	}

	if command.Save {
		if err := d.memory.SaveService(s); err != nil {
			conn.Write([]byte("Error saving new service: " + err.Error() + "\n"))
			return
		}
	}
	startTunnelService(s)
	conn.Write([]byte("New service created and saved successfully\n"))
}

func startTunnelService(s *models.ClientService) {
	if s.TCP.Iniitialized {
		reversetunnel.StartLocalTCPTunnel(s.TCP.Remote, s.TCP.Local)
	}
}
