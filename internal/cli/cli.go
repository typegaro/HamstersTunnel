package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/typegaro/HamstersTunnel/pkg/command"
)

// CLI rappresenta la nostra struttura principale per gestire i comandi
type CLI struct{}

// getSocketPath determina il percorso del socket in base al sistema operativo
func (cli *CLI) getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon`
	}
	return "/tmp/mydaemon.sock"
}

// sendCommand invia un comando al daemon tramite un socket
func (cli *CLI) sendCommand(cmd interface{}) {
	socketPath := cli.getSocketPath()

	conn, err := net.Dial("unix", socketPath)
	if runtime.GOOS == "windows" {
		conn, err = net.Dial("npipe", socketPath)
	}

	if err != nil {
		log.Fatal("Errore nella connessione al daemon:", err)
	}
	defer conn.Close()

	// Codifica il comando come JSON e lo invia al daemon
	json.NewEncoder(conn).Encode(cmd)

	// Legge la risposta dal socket
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	fmt.Println(string(buffer[:n]))
}

// printHelp stampa la guida per l'uso della CLI
func (cli *CLI) printHelp() {
	fmt.Println("Usage: cli <command> [args]")
    fmt.Println("Commands:")
    fmt.Println("  status      : Retrieve the current status.")
    fmt.Println("  set-status  : Set the status with a specified value.")
    fmt.Println("  new         : Create a new service with the provided details.")
}

// status gestisce il comando "status"
func (cli *CLI) status() {
	cmd := command.StatusCommand{Command: "status"}
	cli.sendCommand(cmd)
}

// setStatus gestisce il comando "set-status"
func (cli *CLI) setStatus(value string) {
	cmd := command.SetStatusCommand{
		Command: "set-status",
		Value:   value,
	}
	cli.sendCommand(cmd)
}

// newService gestisce il comando "new"
func (cli *CLI) newService(serviceName, serviceLocalPort, remoteIP,save string,) {
	cmd := command.NewServiceCommand{
		Command:     "new",
		ServiceName: serviceName,
		LocalPort:   serviceLocalPort,
		RemoteIP:    remoteIP,
        Save: save,
	}
	cli.sendCommand(cmd)
}

// Execute gestisce l'esecuzione del comando fornito via CLI
func (cli *CLI) Execute() {
	// Verifica che ci siano abbastanza argomenti
	if len(os.Args) < 2 {
        cli.printHelp()
		os.Exit(1)
	}

	// Switch sui comandi passati dalla riga di comando
	switch os.Args[1] {
	case "status":
		cli.status()
	case "set-status":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli set-status <value>")
			os.Exit(1)
		}
		cli.setStatus(os.Args[2])
	case "new":
		if len(os.Args) < 5 {
			fmt.Println("Usage: cli new <service_name> <service_local_port> <remote_ip> <save_service>")
			os.Exit(1)
		}
		cli.newService(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	default:
		fmt.Println("Unknown command:", os.Args[1])
        cli.printHelp()
		os.Exit(1)
	}
}

