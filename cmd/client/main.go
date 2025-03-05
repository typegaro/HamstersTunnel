package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

func proxyConnection(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer src.Close()
	defer dst.Close()

	fmt.Println("Copying data from", src.RemoteAddr(), "to", dst.RemoteAddr())

	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error copying data from %v to %v: %v\n", src.RemoteAddr(), dst.RemoteAddr(), err)
	}

	_, err = io.Copy(src, dst)
	if err != nil {
		fmt.Printf("Error copying data from %v to %v: %v\n", dst.RemoteAddr(), src.RemoteAddr(), err)
	}
}

func StartTCPProxy(publicPort, privatePort string) error {
	listener, err := net.Listen("tcp", ":"+publicPort)
	if err != nil {
		return fmt.Errorf("failed to start listener on port %s: %w", publicPort, err)
	}
	fmt.Println("TCP proxy listening on port:", publicPort)

	go func() {
		defer listener.Close()
		for {
			clientConn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				return
			}

			serverConn, err := net.Dial("tcp", "127.0.0.1:"+privatePort)
			if err != nil {
				fmt.Println("Error connecting to private port:", err)
				clientConn.Close()
				continue
			}

			var wg sync.WaitGroup
			wg.Add(2)

			go proxyConnection(clientConn, serverConn, &wg)
			go proxyConnection(serverConn, clientConn, &wg)

			wg.Wait()
		}
	}()
	return nil
}

func main() {
	publicPort := "64355"  
	privatePort := "25565" 

	err := StartTCPProxy(publicPort, privatePort)
	if err != nil {
		fmt.Println("Errore nel proxy:", err)
		return
	}

	select {}
}

