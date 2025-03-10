package main 

import (
    "github.com/typegaro/HamstersTunnel/internal/daemon"
)

func main() {
	daemon := &daemon.Daemon{}
	daemon.Init()
}

