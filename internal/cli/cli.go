package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/typegaro/HamstersTunnel/pkg/command"
)

type CLI struct{
    TCP string
    UDP string
    HTTP string
    Ip string
    Name string
    Save bool
}

func (cli *CLI) getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon`
	}
	return "/tmp/mydaemon.sock"
}

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

	json.NewEncoder(conn).Encode(cmd)

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	fmt.Println(string(buffer[:n]))
}

func (cli *CLI) printHelp() {
	fmt.Println("Usage: cli <command> [args]")
    fmt.Println("Commands:")
    fmt.Println("  new         : Create a new service with the provided details.")
}

func (cli *CLI) status() {
	cmd := command.StatusCommand{Command: "status"}
	cli.sendCommand(cmd)
}

func (cli *CLI) setStatus(value string) {
	cmd := command.SetStatusCommand{
		Command: "set-status",
		Value:   value,
	}
	cli.sendCommand(cmd)
}

func (cli *CLI) newService(ip,name,tcp, udp, http string,save bool) {
	cmd := command.NewServiceCommand{
		Command:     "new",
		ServiceName: name,
        TCP:         tcp,
		UDP:         udp,
		HTTP:        http,
		RemoteIP:    ip,
        Save:        save,
	}
	cli.sendCommand(cmd)
}

func (cli *CLI) Execute() {
    cli.newService(cli.Ip, cli.Name, cli.TCP, cli.UDP, cli.HTTP, cli.Save)
}


