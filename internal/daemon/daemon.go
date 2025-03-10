package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
    "bytes"
	"io"
	"net/http"

	"github.com/typegaro/HamstersTunnel/internal/daemon/cache"
	"github.com/typegaro/HamstersTunnel/pkg/command" // Importa il pacchetto command
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type Daemon struct {
	status   string
	listener net.Listener
	cache    cache.FileSystemCache
}

func getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon` // Named Pipe on Windows
	}
	return "/tmp/mydaemon.sock" // Unix socket on Linux/macOS
}

func wakeupTunnel() {
	// TODO: implement me
}

func (d *Daemon) Init() {
	d.cache.Init()
	d.status = "idle"
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

	// Decodifica del comando ricevuto
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

// handleStatus restituisce lo stato attuale del daemon
func (d *Daemon) handleStatus(conn net.Conn) {
	response := map[string]string{"status": d.status}
	json.NewEncoder(conn).Encode(response)
}

// handleSetStatus imposta un nuovo stato per il daemon
func (d *Daemon) handleSetStatus(conn net.Conn, command command.SetStatusCommand) {
	d.status = command.Value
	conn.Write([]byte("Status updated\n"))
}

// handleNewService crea un nuovo servizio
func (d *Daemon) handleNewService(conn net.Conn, command command.NewServiceCommand) {
	// Effettuare una chiamata POST a command.RemoteIP/Service?save=command.Save
	url := fmt.Sprintf("http://%s/Service?save=%v", command.RemoteIP, command.Save)
	
	// Creazione del payload per la richiesta POST in base alla struttura NewServiceReq
	payload := models.NewServiceReq{
		Name:          command.ServiceName,
		TCP:           true,  // Impostato su true, puoi modificarlo in base alle tue necessità
		UDP:           true,  // Impostato su true, puoi modificarlo in base alle tue necessità
		HTTP:          true,  // Impostato su true, puoi modificarlo in base alle tue necessità
		PortBlackList: []string{},  // Lista di porte bloccate, può essere modificata
		PortWitheList: []string{},  // Lista di porte permesse, può essere modificata
		Options:       []string{},  // Eventuali opzioni aggiuntive
	}

	// Codifica il payload in JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		conn.Write([]byte("Error encoding JSON: " + err.Error() + "\n"))
		return
	}

	// Effettua la richiesta POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		conn.Write([]byte("Error making POST request: " + err.Error() + "\n"))
		return
	}
	defer resp.Body.Close()

	// Legge la risposta dalla richiesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		conn.Write([]byte("Error reading response body: " + err.Error() + "\n"))
		return
	}

	// Verifica che la risposta sia di tipo ServiceRes
	var serviceRes models.ServiceRes
	if err := json.Unmarshal(body, &serviceRes); err != nil {
		conn.Write([]byte("Error decoding response body: " + err.Error() + "\n"))
		return
	}

	// Crea un nuovo CachedService usando i dati ricevuti dalla risposta
	cs := &models.CachedService{
		Id:        serviceRes.Id, 
		Name:      command.ServiceName, 
		Ip:        command.RemoteIP,
		HTTP:      serviceRes.HTTP, 
		TCP:       serviceRes.TCP, 
		UDP:       serviceRes.UDP, 
		Active:    true,
	}

	// Salva il nuovo servizio nella cache
	if err := d.cache.SaveService(cs); err != nil {
		conn.Write([]byte("Error saving new service: " + err.Error() + "\n"))
		return
	}

	// Risposta al client
	conn.Write([]byte("New service created and saved successfully\n"))
}
