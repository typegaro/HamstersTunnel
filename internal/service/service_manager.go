package service

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"io"
    "sync"

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
            err := StartTCPProxy(srv.TCP.Public, srv.TCP.Private)
            if err != nil {
                fmt.Printf("Warning: failed to start TCP proxy for service %s: %v\n", srv.Info.Id, err)
                continue // Continua invece di restituire l'errore
            }
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

		err = StartTCPProxy(publicPort, privatePort)
		if err != nil {
			return ps, fmt.Errorf("failed to start TCP proxy: %w", err)
		}

		ps.TCP = &models.PortPair{
			Public:  publicPort,
			Private: privatePort,
		}
	}

	return ps, nil
}

func proxyConnection(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer src.Close()
	defer dst.Close()

	fmt.Println("Copying data from", src.RemoteAddr(), "to", dst.RemoteAddr())

	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error copying data from %v to %v: %v\n", src.RemoteAddr(), dst.RemoteAddr(), err)
	}

	//time.Sleep(10 * time.Millisecond)

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

