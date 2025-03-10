package service

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

// Test function for bidirectional TCP proxy
func TestTCPProxy(t *testing.T) {
    publicPort := "8080"
    privatePort := "9090"

    readyChan := make(chan struct{})

    // Start the TCP proxy in a goroutine
    go func() {
        err := StartTCPProxy(publicPort, privatePort)
        if err != nil {
            t.Errorf("Error starting proxy: %v", err)
        }
    }()

    time.Sleep(1 * time.Second)

    testClient := func(addr string, msg string) (string, error) {
        conn, err := net.Dial("tcp", addr)
        if err != nil {
            return "", fmt.Errorf("error connecting to proxy: %v", err)
        }
        defer conn.Close()

        _, err = conn.Write([]byte(msg))
        if err != nil {
            return "", fmt.Errorf("error sending message: %v", err)
        }

        buffer := make([]byte, len(msg))
        _, err = io.ReadFull(conn, buffer)
        if err != nil {
            return "", fmt.Errorf("error reading response: %v", err)
        }

        return string(buffer), nil
    }

    t.Run("Bidirectional test", func(t *testing.T) {
        go func() {
            listener, err := net.Listen("tcp", ":"+privatePort)
            if err != nil {
                t.Fatalf("Error starting server on private port: %v", err)
            }
            defer listener.Close()

            // Notify that the server is ready
            readyChan <- struct{}{}

            serverConn, err := listener.Accept()
            if err != nil {
                t.Fatalf("Error accepting server connection: %v", err)
            }
            defer serverConn.Close()

            buffer := make([]byte, 1024)
            n, err := serverConn.Read(buffer)
            if err != nil {
                t.Fatalf("Error reading from client: %v", err)
            }

            _, err = serverConn.Write(buffer[:n])
            if err != nil {
                t.Fatalf("Error sending response to client: %v", err)
            }
        }()

        // Wait for the server to be ready
        <-readyChan

        clientMessage := "Hello Server!"
        response, err := testClient("localhost:"+publicPort, clientMessage)
        if err != nil {
            t.Fatalf("Client error: %v", err)
        }

        if response != clientMessage {
            t.Errorf("Message mismatch, expected: %s, got: %s", clientMessage, response)
        }
    })
}

func TestLoadService(t *testing.T) {
    sm := NewServiceManager()
    sm.Init()

    PortBlackList := []string{}
    publicPort, err := findAvailablePort(PortBlackList)
    if err != nil {
        t.Errorf("Error finding available public port: %v", err)
    }

    privatePort, err := findAvailablePort(PortBlackList)
    if err != nil {
        t.Fatalf("Error finding available private port: %v", err)
    }

    srv := &models.PublicService{
        Info: models.ServiceInfo{
            Id:   "f5878e6d-8315-4e05-8c93-f7e959a34580",
            Name: "testname",
        },
        TCP: &models.PortPair{
            Public:  publicPort,
            Private: privatePort,
        },
        UDP:   nil,
        HTTP:  nil,
        Active: true,
    }
    sm.memory.SaveService(srv)

    if _, exists := sm.services[srv.Info.Id]; !exists {
        t.Errorf("Service with id %s not found in sm.services", srv.Info.Id)
    }
}
