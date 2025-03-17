package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/typegaro/HamstersTunnel/pkg/command"
)

type CLI struct{}

func (cli *CLI) getSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\mydaemon`
	}
	return "/tmp/mydaemon.sock"
}

func (cli *CLI) sendCommand(cmd interface{}) {
	socketPath := cli.getSocketPath()

	var conn net.Conn
	var err error
	if runtime.GOOS == "windows" {
		conn, err = net.Dial("npipe", socketPath)
	} else {
		conn, err = net.Dial("unix", socketPath)
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

func (cli *CLI) NewService(ip, name, tcp, udp, http string, save bool) {
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

// TODO: Implement this on deamon side
func (cli *CLI) ListService(inactive bool) {
	cmd := command.ListCommand{Command: "status", Inactive: inactive}
	cli.sendCommand(cmd)
}

// TODO: Implement this on deamon side
func (cli *CLI) StopService(id string, remote bool) {
	cmd := command.ServiceCommand{Command: "stop", Id: id, Remote: remote}
	cli.sendCommand(cmd)
}

// TODO: Implement this on deamon side
func (cli *CLI) RemoveService(id string, remote bool) {
	cmd := command.ServiceCommand{Command: "remove", Id: id, Remote: remote}
	cli.sendCommand(cmd)
}
