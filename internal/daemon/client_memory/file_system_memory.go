package client_memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/typegaro/HamstersTunnel/pkg/models/service"
	"github.com/typegaro/HamstersTunnel/pkg/utility"
)

type FileSystemMemory struct {
	storagePath string
	services    map[string]*models.ClientService
}

// Initializes the storage directory
func (fs *FileSystemMemory) Init() {
	baseDir, _ := os.Getwd()
	fs.storagePath = filepath.Join(baseDir, ".config/data")

	if err := os.MkdirAll(fs.storagePath, os.ModePerm); err != nil {
		panic("failed to create directory: " + err.Error())
	}

	fs.services = make(map[string]*models.ClientService)
	if srvs, err := fs.getServicesFromFile(); err != nil {
		for _, srv := range srvs {
			fs.services[srv.Id] = srv
		}
	}
}

func (fd *FileSystemMemory) AddService(srv *models.ClientService) error {
	if !fd.IsService(srv.Id) {
		fd.services[srv.Id] = srv
		fd.saveServiceOnFile(srv)
	} else {
		return fmt.Errorf("Error: %s already exist", srv.Id)
	}
	return nil
}

func (fd *FileSystemMemory) RemoveService(id string) error {
	if !fd.IsService(id) {
		fd.services[id] = nil
		fd.deleteFile(id)
	} else {
		return fmt.Errorf("Error: %s don't exist", id)
	}
	return nil
}

func (fd *FileSystemMemory) EditService(srv *models.ClientService) error {
	if fd.IsService(srv.Id) {
		fd.services[srv.Id] = srv
		if err := fd.saveServiceOnFile(srv); err != nil {
			return err
		}
	}
	return nil
}

func (fd *FileSystemMemory) GetServices() []*models.ClientService {
	return utitlity.MapGetValues(fd.services)
}

func (fd *FileSystemMemory) GetService(id string) *models.ClientService {
	return fd.services[id]
}

// Retrieves all active services
func (fs *FileSystemMemory) getActiveServicesFromFile() ([]*models.ClientService, error) {
	files, err := os.ReadDir(fs.storagePath)
	if err != nil {
		return nil, err
	}

	var services []*models.ClientService

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		serviceFilePath := filepath.Join(fs.storagePath, file.Name())
		data, err := os.ReadFile(serviceFilePath)
		if err != nil {
			return nil, err
		}

		var srv models.ClientService
		if err := json.Unmarshal(data, &srv); err != nil {
			return nil, err
		}
		if srv.Active {
			services = append(services, &srv)
		}
	}

	return services, nil
}

func (fs *FileSystemMemory) getServicesFromFile() ([]*models.ClientService, error) {
	files, err := os.ReadDir(fs.storagePath)
	if err != nil {
		return nil, err
	}

	var services []*models.ClientService

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		serviceFilePath := filepath.Join(fs.storagePath, file.Name())
		data, err := os.ReadFile(serviceFilePath)
		if err != nil {
			return nil, err
		}

		var srv models.ClientService
		if err := json.Unmarshal(data, &srv); err != nil {
			return nil, err
		}
		services = append(services, &srv)
	}
	return services, nil
}

func (fs *FileSystemMemory) saveServiceOnFile(srv *models.ClientService) error {
	serviceFilePath := filepath.Join(fs.storagePath, srv.Id+".json")

	if _, err := os.Stat(fs.storagePath); os.IsNotExist(err) {
		return fmt.Errorf("destination folder does not exist: %v", fs.storagePath)
	}

	if data, err := json.MarshalIndent(srv, "", "  "); err != nil {
		return err
	} else {
		if err := os.WriteFile(serviceFilePath, data, 0644); err != nil {
			return fmt.Errorf("error saving service: %w", err)
		}
	}

	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s was not created", serviceFilePath)
	}

	return nil
}

// Retrieves a service by ID
func (fs *FileSystemMemory) getServiceFromFile(id string) (*models.ClientService, error) {
	serviceFilePath := filepath.Join(fs.storagePath, id+".json")

	if file, err := os.ReadFile(serviceFilePath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("service not found")
		}
		return nil, err
	} else {
		var srv models.ClientService
		if err := json.Unmarshal(file, &srv); err != nil {
			return nil, err
		}
		return &srv, nil
	}
}

// Deletes a service by ID
func (fs *FileSystemMemory) deleteFile(id string) error {

	serviceFilePath := filepath.Join(fs.storagePath, id+".json")

	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		return errors.New("service not found")
	}

	return os.Remove(serviceFilePath)
}

// Checks if a service exists
func (fs *FileSystemMemory) IsService(id string) bool {
	serviceFilePath := filepath.Join(fs.storagePath, id+".json")

	_, err := os.Stat(serviceFilePath)
	return !os.IsNotExist(err)
}
