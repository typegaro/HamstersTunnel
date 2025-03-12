package main 

import (
    "github.com/typegaro/HamstersTunnel/internal/daemon"
)

func main() {
	daemon := daemon.NewDaemon()
	daemon.Init()
}

