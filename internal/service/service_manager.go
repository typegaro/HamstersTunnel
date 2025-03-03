package service

import (
    "github.com/google/uuid"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
    "fmt"

)

type ServiceManager struct{
    usedPorts map[string]string
    services map[string]*models.PublicService
}

func NewServiceManager() *ServiceManager {
    return &ServiceManager{
        usedPorts: make(map[string]string),
        services: make(map[string]*models.PublicService),
    }
}

func (sm *ServiceManager) addService(ps *models.PublicService) error{
    if _, exists := sm.services[ps.Info.Id]; exists {
        return fmt.Errorf("service with id %s already exists", ps.Info.Id)
    }
    sm.services[ps.Info.Id] = ps
    return nil
}

func (sm *ServiceManager) removeService(ps *models.PublicService) error{
    if _, exists := sm.services[ps.Info.Id]; exists {
        return fmt.Errorf("service with id %s not found", ps.Info.Id)
    }
    delete(sm.services, ps.Info.Id)
    return nil
}

func (sm *ServiceManager) isUsedPort(port string) (bool,string) {
    if value, exists := sm.usedPorts[port]; exists {
        return true,value
    }
    return false,""
}

func (sm *ServiceManager) addUsedPort(port string,id string) error{
    if value,user :=sm.isUsedPort(port); !value{
        return fmt.Errorf("Port used by %s", user)
    }
    sm.usedPorts[port]= id   
    return nil
}

// GeneratePublicService creates a PublicService from a NewServiceReq
func GeneratePublicService(req models.NewServiceReq) (models.PublicService, error) {
    ps := models.PublicService{}
    var err error

    if req.TCP {
        ps.TCP,err = makePair("tcp")
        if err != nil {
            return ps, fmt.Errorf("error generating TCP port pair: %w", err)
        }
    }

    if req.UDP {
        ps.UDP,err = makePair("udp")
        if err != nil {
            return ps, fmt.Errorf("error generating UDP port pair: %w", err)
        }
    }

    if req.HTTP {
        ps.HTTP,err = makePair("http")
        if err != nil {
            return ps, fmt.Errorf("error generating HTTP port pair: %w", err)
        }
    }

    // Set a default service info for the public service
    ps.Info = models.ServiceInfo{
        Id:   uuid.New().String(), 
        Name: req.Name,
    }

    return ps, nil
}

// makePair generates a port pair for a given service (HTTP, TCP, UDP)
func makePair(serviceType string) (models.PortPair, error) {
    // TODO: Generate real port pair logic (either generate random or specific logic)
    // Just as an example, let's assume we generate some default ports based on service type
    var pair models.PortPair

    switch serviceType {
    case "http":
        pair = models.PortPair{
            External: "80",
            Internal: "8080",
        }
    case "tcp":
        pair = models.PortPair{
            External: "2023",
            Internal: "3000",
        }
    case "udp":
        pair = models.PortPair{
            External: "4021",
            Internal: "4000",
        }
    default:
        return models.PortPair{}, fmt.Errorf("unsupported service type: %s", serviceType)
    }

    return pair, nil
}
