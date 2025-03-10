package server 

import (
	"github.com/typegaro/HamstersTunnel/internal/server/service"
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

func (s *Server) Init() error {
    if err:=s.ServiceManager.Init(); err!= nil{
        return err
    }
    return nil
}
