package main

import(
    "flag"
    "fmt"
    "github.com/typegaro/HamstersTunnel/internal/cli"
)

func main() {
    tcp := flag.String("tcp", "", "Set local service TCP port") 
    udp := flag.String("udp", "", "Set local service UDP port") 
    http := flag.String("http", "", "Set local service HTTP port") 
    ip := flag.String("ip", "", "Server ip and port (ip:port)") 
    name := flag.String("name", "", "Service name") 
    save := flag.Bool("s", false, "Save the service on remote for quick restart")
    flag.Parse()

    fmt.Println("Name:", *name)
    fmt.Println("IP:", *ip)
    fmt.Println("Save:", *save)

	cli := &cli.CLI{
        TCP: *tcp,
        UDP: *udp,
        HTTP: *http,
        Ip: *ip,
        Name: *name,
        Save: *save,
    }
	cli.Execute()
}

