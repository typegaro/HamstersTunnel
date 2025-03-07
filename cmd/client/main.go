package main

import (
    "log"
    "net"
    "io"
    "errors"
    "strings"
    "os"
    "fmt"
)

// Check the reason for the connection closure
func checkConnectionClosed(err error) string {
    if err == nil {
        return "Connection is active"
    }

    if errors.Is(err, io.EOF) {
        return "Normal closure (FIN received)"
    }
    if strings.Contains(err.Error(), "connection reset by peer") {
        return "Forced closure (RST received)"
    }
    if strings.Contains(err.Error(), "i/o timeout") {
        return "Closure due to inactivity timeout"
    }
    if strings.Contains(err.Error(), "broken pipe") {
        return "Closure due to network disruption"
    }

    return "Unknown error: " + err.Error()
}

func forwardData(src, dst net.Conn) {
    defer src.Close()
    defer dst.Close()

    go func() {
        _, err := io.Copy(dst, src)
        if err != nil {
            if !errors.Is(err, net.ErrClosed) {
                log.Printf("Error forwarding data from src to dst: %v (%s)", err, checkConnectionClosed(err))
            }
        } else {
            log.Println("Connection closed normally from src to dst")
        }
    }()

    _, err := io.Copy(src, dst)
    if err != nil {
        if !errors.Is(err, net.ErrClosed) {
            log.Printf("Error forwarding data from dst to src: %v (%s)", err, checkConnectionClosed(err))
        }
    } else {
        log.Println("Connection closed normally from dst to src")
    }
}

func startProxy(remotePort, servicePort string) {
    remoteConn, err := net.Dial("tcp", remotePort)
    if err != nil {
        log.Fatalf("Unable to connect to remote proxy on port %s: %v", remotePort, err)
    }
    defer remoteConn.Close()
    log.Printf("Connected to remote proxy: %s", remotePort)

    // Connect to the local service
    serviceConn, err := net.Dial("tcp", servicePort)
    if err != nil {
        log.Fatalf("Error connecting to the local service on port %s: %v", servicePort, err)
    }
    defer serviceConn.Close()
    log.Printf("Connected to local service: %s", servicePort)

    // Forward data in both directions
    go forwardData(remoteConn, serviceConn)
    go forwardData(serviceConn, remoteConn)

    // Block the main function to keep connections active
    select {}
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run main.go <remotePort> <servicePort>")
        os.Exit(1)
    }

    remotePort := os.Args[1]
    servicePort := os.Args[2]

    startProxy(remotePort, servicePort)
}

