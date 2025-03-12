package memory 

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
    "fmt"

	"github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type FileSystemMemory struct {
	storagePath string
	mutex       sync.Mutex
}

// Initializes the storage directory
func (fs *FileSystemMemory) Init() {
	baseDir, _ := os.Getwd()
	fs.storagePath = filepath.Join(baseDir, "data")

	if err := os.MkdirAll(fs.storagePath, os.ModePerm); err != nil {
		panic("failed to create directory: " + err.Error())
	}
}

// Saves a service to the file system
func (fs *FileSystemMemory) SaveService(srv *models.Service) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

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
func (fs *FileSystemMemory) GetService(id string) (*models.Service, error) {
	serviceFilePath := filepath.Join(fs.storagePath, id+".json")

	if file, err := os.ReadFile(serviceFilePath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}else{
        var srv models.Service
	    if err := json.Unmarshal(file, &srv); err != nil {
	    	return nil, err
	    }
	    return &srv, nil
    }
}

// Retrieves all active services
func (fs *FileSystemMemory) GetActiveServices() ([]*models.Service, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	files, err := os.ReadDir(fs.storagePath)
	if err != nil {
		return nil, err
	}

	var services []*models.Service

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

		var srv models.Service
		if err := json.Unmarshal(data, &srv); err != nil {
			return nil, err
		}

		services = append(services, &srv)
	}

	return services, nil
}
    
// Deletes a service by ID
func (fs *FileSystemMemory) DeleteService(id string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

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

