package reversetunnel

import (
	"io"
	"log"
	"net"
)

func forwardData(clientConn, userConn net.Conn) {
	go func() {
		_, err := io.Copy(userConn, clientConn)
		if err != nil {
			if err.Error() == "use of closed network connection" {
				log.Println("Connection closed by client while forwarding data to user.")
			} else {
				log.Printf("Error forwarding data from client to user: %v", err)
			}
		} else {
			log.Println("Data forwarding from client to user completed.")
		}
	}()

	_, err := io.Copy(clientConn, userConn)
	if err != nil {
		if err.Error() == "use of closed network connection" {
			log.Println("Connection closed by user while forwarding data to client.")
		} else {
			log.Printf("Error forwarding data from user to client: %v", err)
		}
	} else {
		log.Println("Data forwarding from user to client completed.")
	}
}

func StartRemoteTCPTunnel(clientPort, userPort string) error {
	clientListener, err := net.Listen("tcp", ":"+clientPort)
	if err != nil {
		log.Fatalf("Unable to start listener on public port %s: %v", clientPort, err)
	}
	//defer clientListener.Close()

	userListener, err := net.Listen("tcp", ":"+userPort)
	if err != nil {
		log.Fatalf("Unable to start listener on user port %s: %v", userPort, err)
	}
	defer userListener.Close()

	log.Printf("Remote proxy listening on ports %s (client) and %s (user)", clientPort, userPort)

	clientConn, err := clientListener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection from client: %v", err)
	}
	log.Println("Connection established with client.")

	for {
		userConn, err := userListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from user: %v", err)
			continue
		}

		log.Println("Connection established with user.")
		go forwardData(clientConn, userConn)
	}
}
