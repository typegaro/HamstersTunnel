package service

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"io"
	"sync"
	"log"

	"github.com/google/uuid"
	"github.com/typegaro/HamstersTunnel/internal/memory"
	"github.com/typegaro/HamstersTunnel/pkg/interfaces"
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type ServiceManager struct {
	usedPorts map[string]string
	services  map[string]*models.PublicService
	memory    interfaces.AbstractMemory
	mutex     sync.Mutex
}

func NewServiceManager() *ServiceManager {
	rand.Seed(time.Now().UnixNano())

	return &ServiceManager{
		usedPorts: make(map[string]string),
		services:  make(map[string]*models.PublicService),
		memory:    &memory.FileSystemMemory{},
	}
}

func (sm *ServiceManager) Init() error {
	sm.memory.Init()
	if err := sm.loadService(); err != nil {
		return fmt.Errorf("failed to initialize service manager: %w", err)
	}
	return nil
}

func (sm *ServiceManager) loadService() error {
	srvs, err := sm.memory.GetActiveServices()
	if err != nil {
		return err
	}

	for _, srv := range srvs {
		if srv.TCP != nil {
			go StartTCPProxy(srv.TCP.Public, srv.TCP.Private)
		}
		sm.addService(srv)
	}
	return nil
}

func (sm *ServiceManager) addService(ps *models.PublicService) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.services[ps.Info.Id]; exists {
		return fmt.Errorf("service with id %s already exists", ps.Info.Id)
	}
	sm.services[ps.Info.Id] = ps
	return nil
}

func (sm *ServiceManager) removeService(ps *models.PublicService) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.services[ps.Info.Id]; !exists {
		return fmt.Errorf("service with id %s not found", ps.Info.Id)
	}
	delete(sm.services, ps.Info.Id)
	return nil
}

func (sm *ServiceManager) isUsedPort(port string) (string, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	value, exists := sm.usedPorts[port]
	return value, exists
}

func (sm *ServiceManager) addUsedPort(port string, id string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if user, exists := sm.usedPorts[port]; exists {
		return fmt.Errorf("port %s is already used by service %s", port, user)
	}
	sm.usedPorts[port] = id
	return nil
}

func GeneratePublicService(req models.NewServiceReq) (models.PublicService, error) {
	ps := models.PublicService{
		Info: models.ServiceInfo{
			Id:   uuid.New().String(),
			Name: req.Name,
		},
		Active: true,
	}

	if req.TCP {
		publicPort, err := findAvailablePort(req.PortBlackList)
		if err != nil {
			return ps, fmt.Errorf("failed to find an available public port: %w", err)
		}

		privatePort, err := findAvailablePort(req.PortBlackList)
		if err != nil {
			return ps, fmt.Errorf("failed to find an available private port: %w", err)
		}

		go StartTCPProxy(publicPort, privatePort)

		ps.TCP = &models.PortPair{
			Public:  publicPort,
			Private: privatePort,
		}
	}

	return ps, nil
}

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

func StartTCPProxy(publicPort, proxyPort string) error {
	publicListener, err := net.Listen("tcp", ":"+publicPort)
	if err != nil {
		log.Fatalf("Unable to start listener on public port %s: %v", publicPort, err)
	}
	defer publicListener.Close()

	proxyListener, err := net.Listen("tcp", ":"+proxyPort)
	if err != nil {
		log.Fatalf("Unable to start listener on proxy port %s: %v", proxyPort, err)
	}
	defer proxyListener.Close()

	log.Printf("Remote proxy listening on ports %s (client) and %s (proxy)", publicPort, proxyPort)

	proxyConn, err := proxyListener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection from local proxy: %v", err)
	}
	log.Println("Connection established with local proxy.")

	for {
		clientConn, err := publicListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %v", err)
			continue
		}

		log.Println("Connection established with client.")
		go forwardData(clientConn, proxyConn)
		go forwardData(proxyConn, clientConn)
	}
}

func findAvailablePort(blacklist []string) (string, error) {
	for i := 0; i < 100; i++ {
		port := generateRandomPort()
		portStr := strconv.Itoa(port)

		if !contains(blacklist, portStr) && isPortAvailable(port) {
			return portStr, nil
		}
	}
	return "", fmt.Errorf("unable to find an available port")
}

func generateRandomPort() int {
	return rand.Intn(64512) + 1024
}

func isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func contains(list []string, port string) bool {
	for _, item := range list {
		if item == port {
			return true
		}
	}
	return false
}

