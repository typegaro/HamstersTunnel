package reversetunnel

import (
	"net"
	"io"
	"log"
)

func forwardData(src, dst net.Conn) {
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			if err.Error() == "use of closed network connection" {
				log.Println("Connection closed by src while forwarding data to dst.")
			} else {
				log.Printf("Error forwarding data from src to dst: %v", err)
			}
		} else {
			log.Println("Data forwarding from src to dst completed.")
		}
	}()

	_, err := io.Copy(src, dst)
	if err != nil {
		if err.Error() == "use of closed network connection" {
			log.Println("Connection closed by dst while forwarding data to src.")
		} else {
			log.Printf("Error forwarding data from dst to src: %v", err)
		}
	} else {
		log.Println("Data forwarding from dst to src completed.")
	}
}

func StartRemoteTCPTunnel(clientPort, proxyPort string) error {
	clientListener, err := net.Listen("tcp", ":"+clientPort)
	if err != nil {
		log.Fatalf("Unable to start listener on public port %s: %v", clientPort, err)
	}
	defer clientListener.Close()

	proxyListener, err := net.Listen("tcp", ":"+proxyPort)
	if err != nil {
		log.Fatalf("Unable to start listener on proxy port %s: %v", proxyPort, err)
	}
	defer proxyListener.Close()

	log.Printf("Remote proxy listening on ports %s (client) and %s (proxy)", clientPort, proxyPort)

	proxyConn, err := proxyListener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection from local proxy: %v", err)
	}
	log.Println("Connection established with local proxy.")

	for {
		clientConn, err := clientListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %v", err)
			continue
		}

		log.Println("Connection established with client.")
		go forwardData(clientConn, proxyConn)
		go forwardData(proxyConn, clientConn)
	}
}
