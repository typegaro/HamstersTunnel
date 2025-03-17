package service

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/typegaro/HamstersTunnel/internal/server/server_memory"
	"github.com/typegaro/HamstersTunnel/pkg/interfaces"
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
	"github.com/typegaro/HamstersTunnel/pkg/reversetunnel"
)

type ServiceManager struct {
	usedPorts map[string]string
	services  map[string]*models.ServerService
	memory    interfaces.ServerMemory
}

func NewServiceManager() *ServiceManager {
	rand.Seed(time.Now().UnixNano())

	return &ServiceManager{
		usedPorts: make(map[string]string),
		services:  make(map[string]*models.ServerService),
		memory:    &server_memory.FileSystemMemory{},
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
			go reversetunnel.StartRemoteTCPTunnel(srv.TCP.Proxy, srv.TCP.Client)
		}
		sm.addService(srv)
	}
	return nil
}

func (sm *ServiceManager) addService(ps *models.ServerService) error {

	if _, exists := sm.services[ps.Id]; exists {
		return fmt.Errorf("service with id %s already exists", ps.Id)
	}
	sm.services[ps.Id] = ps
	return nil
}

func (sm *ServiceManager) removeService(ps *models.ServerService) error {

	if _, exists := sm.services[ps.Id]; !exists {
		return fmt.Errorf("service with id %s not found", ps.Id)
	}
	delete(sm.services, ps.Id)
	return nil
}

func (sm *ServiceManager) isUsedPort(port string) (string, bool) {

	value, exists := sm.usedPorts[port]
	return value, exists
}

func (sm *ServiceManager) addUsedPort(port string, id string) error {

	if user, exists := sm.usedPorts[port]; exists {
		return fmt.Errorf("port %s is already used by service %s", port, user)
	}
	sm.usedPorts[port] = id
	return nil
}

func GeneratePublicService(req models.NewServiceReq) (models.ServerService, error) {
	ps := models.ServerService{
		Id:      uuid.New().String(),
		Name:    req.Name,
		Active:  true,
		Options: []string{},
	}

	if req.TCP {
		proxyPort, err := findAvailablePort(req.PortBlackList)
		if err != nil {
			return ps, fmt.Errorf("failed to find an available public port: %w", err)
		}

		clientPort, err := findAvailablePort(req.PortBlackList)
		if err != nil {
			return ps, fmt.Errorf("failed to find an available private port: %w", err)
		}

		go reversetunnel.StartRemoteTCPTunnel(proxyPort, clientPort)

		ps.TCP = &models.ServerPortPair{
			Proxy:  proxyPort,
			Client: clientPort,
		}
	}
	return ps, nil
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
