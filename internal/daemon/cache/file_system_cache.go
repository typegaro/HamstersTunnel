package cache

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
)

type FileSystemCache struct {
    memory map[string]*models.CachedService
}

func (cache *FileSystemCache) Init() {
    cache.memory = make(map[string]*models.CachedService)

    if err := loadFromFile(&cache.memory); err != nil {
        fmt.Println("Warning: Unable to load data from file:", err)
    }
}

func loadFromFile(memory *map[string]*models.CachedService) error {
    file, err := os.Open("cache.json")
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(memory); err != nil {
        return err
    }

    return nil
}

func saveOnFile(memory map[string]*models.CachedService) error {
    file, err := os.Create("cache.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") 
    if err := encoder.Encode(memory); err != nil {
        return err
    }

    return nil
}

func (cache *FileSystemCache) SaveService(srv *models.CachedService) error {
    if srv == nil {
        return errors.New("service cannot be nil")
    }

    if srv.Id == "" {
        return errors.New("service must have a valid ID")
    }

    cache.memory[srv.Id] = srv

    if err := saveOnFile(cache.memory); err != nil {
        return fmt.Errorf("failed to save service to file: %v", err)
    }

    return nil
}

func (cache *FileSystemCache) GetService(id string) (*models.CachedService, error) {
    if service, found := cache.memory[id]; found {
        return service, nil
    }
    return nil, errors.New("service not found")
}

func (cache *FileSystemCache) GetActiveServices() ([]*models.CachedService, error) {
    var activeServices []*models.CachedService
    for _, service := range cache.memory {
        if service.Active {
            activeServices = append(activeServices, service)
        }
    }
    return activeServices, nil
}

func (cache *FileSystemCache) DeleteService(id string) error {
    if _, found := cache.memory[id]; found {
        delete(cache.memory, id)

        if err := saveOnFile(cache.memory); err != nil {
            return fmt.Errorf("failed to save updated data to file: %v", err)
        }
        return nil
    }
    return errors.New("service not found")
}

func (cache *FileSystemCache) IsService(id string) bool {
    _, found := cache.memory[id]
    return found
}

