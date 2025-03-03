package server 

import (
	"github.com/typegaro/HamstersTunnel/internal/service"

)

type Server struct{
    ServiceManager *service.ServiceManager
}

func NewServer() *Server{
    return &Server{
        ServiceManager: service.NewServiceManager(),
    }
}
