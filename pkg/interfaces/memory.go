package interfaces

import (
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type ServerMemory interface {
	Init()

	AddService(srv *models.ServerService) error

	RemoveService(id string) error

	EditService(srv *models.ServerService) error

	GetServices() []*models.ServerService

	GetService(id string) *models.ServerService

	IsService(id string) bool
}

type ClientMemory interface {
	Init()

	AddService(srv *models.ClientService) error

	RemoveService(id string) error

	EditService(srv *models.ClientService) error

	GetServices() []*models.ClientService

	GetService(id string) *models.ClientService

	IsService(id string) bool
}
