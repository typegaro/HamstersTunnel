package interfaces

import (
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type AbstractMemory interface {
    // Init method
    Init()

    // SaveService stores a PublicService instance in memory.
    // Returns an error if the service cannot be saved.
    SaveService(srv *models.PublicService) error

    // GetService retrieves a PublicService instance by its ID.
    // Returns the service if found, but consider handling cases where the ID does not exist.
    GetService(id string) (*models.PublicService, error)

    // GetActiveService search and returns the active service
    GetActiveServices() ([]*models.PublicService, error) 

    // DeleteService removes a PublicService instance by its ID.
    // Returns an error if the service cannot be deleted.
    DeleteService(id string) error

    // IsService checks whether a service with the given ID exists.
    IsService(id string) bool
}
