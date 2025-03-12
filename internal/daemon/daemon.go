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

	"github.com/typegaro/HamstersTunnel/internal/shared/memory"
	"github.com/typegaro/HamstersTunnel/pkg/command"
	"github.com/typegaro/HamstersTunnel/pkg/interfaces"
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
	"github.com/typegaro/HamstersTunnel/pkg/reversetunnel"
)

type Daemon struct {
	status   string
	listener net.Listener
	memory interfaces.Memory
}

func NewDaemon() *Daemon{
    return &Daemon{
        status: "Runnig",
        memory: &memory.FileSystemMemory{},
    }
}

func getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon` // Named Pipe on Windows
	}
	return "/tmp/mydaemon.sock" // Unix socket on Linux/macOS
}

func (d *Daemon) Init() {
	d.memory.Init()
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
	defer d.listener.Close()

	fmt.Println("Daemon in ascolto su", socketPath)

	for {
		conn, err := d.listener.Accept()
		if err != nil {
			log.Println("Errore nella connessione:", err)
			continue
		}
		go d.handleConnection(conn)
	}
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
	case "set-status":
		var setStatusCommand command.SetStatusCommand
		if err := json.Unmarshal(buffer[:n], &setStatusCommand); err != nil {
			conn.Write([]byte("Errore nel parsing del comando 'set-status'\n"))
			return
		}
		d.handleSetStatus(conn, setStatusCommand)
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

func (d *Daemon) handleSetStatus(conn net.Conn, command command.SetStatusCommand) {
	d.status = command.Value
	conn.Write([]byte("Status updated\n"))
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

	s := &models.Service{
		Id:        serviceRes.Id, 
        Name:      command.ServiceName+":"+command.RemoteIP, 
		TCP:       &models.PortPair{Proxy: serviceRes.TCP,Client: command.TCP}, 
		UDP:       &models.PortPair{Proxy: serviceRes.UDP,Client: command.UDP}, 
        HTTP:      &models.PortPair{Proxy: serviceRes.HTTP,Client: command.HTTP}, 
		Active:    true,
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

func startTunnelService(s *models.Service){
    if s.TCP.Client != ""{
        reversetunnel.StartLocalTCPTunnel(s.TCP.Proxy,s.TCP.Proxy) 
    }
}
