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
	"github.com/typegaro/HamstersTunnel/pkg/utility"
)

type Daemon struct {
	status   string
	listener net.Listener
	memory   interfaces.ClientMemory
}

func NewDaemon() *Daemon {
	return &Daemon{
		status: "Running",
		memory: &client_memory.FileSystemMemory{},
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
			log.Println("Connection error:", err)
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
		log.Fatal("Error opening socket:", err)
	}

	fmt.Println("Daemon listening on", socketPath)
}

func (d *Daemon) WakeUpServices() error {
	for _, srv := range d.memory.GetServices() {
		if srv.TCP != nil {
			go reversetunnel.StartLocalTCPTunnel(srv.TCP.Remote, srv.TCP.Local)
		}
	}
	return nil
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	var request command.Command
	err = json.Unmarshal(buffer[:n], &request)
	if err != nil {
		fmt.Println("JSON parsing error:", err)
		return
	}

	switch request.Command {
	case "ls":
		var listCommand command.ListCommand
		if err := json.Unmarshal(buffer[:n], &listCommand); err != nil {
			conn.Write([]byte("Error parsing the 'ls' command\n"))
			return
		}
		d.handleList(conn, listCommand)
	case "stop":
		var stopCommand command.ServiceCommand
		if err := json.Unmarshal(buffer[:n], &stopCommand); err != nil {
			conn.Write([]byte("Error parsing the 'stop' command\n"))
			return
		}
		d.handleStop(conn, stopCommand)
	case "rm":
		var removeCommand command.ServiceCommand
		if err := json.Unmarshal(buffer[:n], &removeCommand); err != nil {
			conn.Write([]byte("Error parsing the 'stop' command\n"))
			return
		}
		d.handleRemove(conn, removeCommand)
	case "start":
		var startCommand command.ServiceCommand
		if err := json.Unmarshal(buffer[:n], &startCommand); err != nil {
			conn.Write([]byte("Error parsing the 'start' command\n"))
			return
		}
		d.handleStart(conn, startCommand)
	case "new":
		var newServiceCommand command.NewServiceCommand
		if err := json.Unmarshal(buffer[:n], &newServiceCommand); err != nil {
			conn.Write([]byte("Error parsing the 'new' command\n"))
			return
		}
		d.handleNewService(conn, newServiceCommand)
	default:
		conn.Write([]byte("Unknown command\n"))
	}
}

func sendHTTPRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	return resp, nil
}

func (d *Daemon) handleStart(conn net.Conn, command command.ServiceCommand) {
	url := fmt.Sprintf(
		"http://%s/service/%s/start",
		d.memory.GetService(command.Id).Ip,
		command.Id,
	)

	resp, err := sendHTTPRequest("POST", url, nil)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	defer resp.Body.Close()
	srv := d.memory.GetService(command.Id)
	srv.Active = true
	if err := d.memory.EditService(srv); err != nil {
		conn.Write([]byte("Error removing service: " + err.Error() + "\n"))
		return
	}

	conn.Write([]byte("Service " + command.Id + " removed \n"))
}

func (d *Daemon) handleRemove(conn net.Conn, command command.ServiceCommand) {
	url := fmt.Sprintf(
		"http://%s/service/%s",
		d.memory.GetService(command.Id).Ip,
		command.Id,
	)

	resp, err := sendHTTPRequest("DELETE", url, nil)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	defer resp.Body.Close()

	if err := d.memory.RemoveService(command.Id); err != nil {
		conn.Write([]byte("Error removing service: " + err.Error() + "\n"))
		return
	}

	conn.Write([]byte("Service " + command.Id + " removed \n"))
}

func (d *Daemon) handleStop(conn net.Conn, command command.ServiceCommand) {
	//Stop remote service
	url := fmt.Sprintf("http://%s/service/%s/stop", command.Remote, command.Id)
	resp, err := sendHTTPRequest("PUT", url, nil)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	//Stop local service
	srv := d.memory.GetService(command.Id)
	srv.Active = false
	defer resp.Body.Close()
	if err := d.memory.EditService(srv); err != nil {
		conn.Write([]byte("Error stoping Service: " + err.Error() + "\n"))
	}
	conn.Write([]byte("Service " + command.Id + " Stopped\n"))
}

func (d *Daemon) handleList(conn net.Conn, command command.ListCommand) {
	var output bytes.Buffer
	output.WriteString("ID\tNAME\tIP\tSTATUS\tTCP\tUDP\tHTTP\n\n")

	for _, s := range d.memory.GetServices() {
		if s.Active || command.Inactive {
			output.WriteString(fmt.Sprintf("%s\t%s\t%s\t%t\t%s\t%s\t%s\n",
				s.Id,
				s.Name,
				s.Ip,
				s.Active,
				utitlity.Ternary(s.TCP != nil, s.TCP.Local+"->"+s.TCP.Remote, "N/A"),
				utitlity.Ternary(s.UDP != nil, s.UDP.Local+"->"+s.HTTP.Remote, "N/A"),
				utitlity.Ternary(s.HTTP != nil, s.HTTP.Local+"->"+s.HTTP.Remote, "N/A"),
			))
		}
	}

	conn.Write(output.Bytes())
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

	var serviceRes models.NewServiceRes
	if err := json.Unmarshal(body, &serviceRes); err != nil {
		conn.Write([]byte("Error decoding response body: " + err.Error() + "\n"))
		return
	}

	s := &models.ClientService{
		Id:   serviceRes.Id,
		Name: command.ServiceName,
		Ip:   command.RemoteIP,
		TCP: utitlity.Ternary(
			serviceRes.TCP != "",
			&models.ClientPortPair{Remote: serviceRes.TCP, Local: command.TCP},
			nil,
		),
		UDP: utitlity.Ternary(
			serviceRes.UDP != "",
			&models.ClientPortPair{Remote: serviceRes.UDP, Local: command.UDP},
			nil,
		),
		HTTP: utitlity.Ternary(
			serviceRes.HTTP != "",
			&models.ClientPortPair{Remote: serviceRes.HTTP, Local: command.HTTP},
			nil,
		),
		Active: true,
	}

	if command.Save {
		if err := d.memory.AddService(s); err != nil {
			conn.Write([]byte("Error saving new service: " + err.Error() + "\n"))
			return
		}
	}
	go startTunnelService(s)

	response, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		conn.Write([]byte("Error encoding response: " + err.Error() + "\n"))
		return
	}
	conn.Write([]byte("New service created and saved successfully\n"))
	conn.Write(response)
}

func startTunnelService(s *models.ClientService) {
	if s.TCP != nil {
		reversetunnel.StartLocalTCPTunnel(s.TCP.Remote, s.TCP.Local)
	}
	//TODO: Implement UDP and HTTP
}
