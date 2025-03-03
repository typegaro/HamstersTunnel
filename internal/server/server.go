package server 

import (
	"github.com/typegaro/HamstersTunnel/internal/service"

)

// Server represents the main application server
type Server struct {
    ServiceManager *service.ServiceManager
}

// NewServer creates and returns a new Server instance
func NewServer() *Server {
    return &Server{
        ServiceManager: service.NewServiceManager(),
    }
}
