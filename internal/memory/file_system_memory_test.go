package memory 

import (
	"github.com/typegaro/HamstersTunnel/pkg/models/service"
	"os"
	"testing"
)

func setupTest(t *testing.T) *FileSystemMemory {
	// Initializes the file system for testing, using a separate test folder
	fs := &FileSystemMemory{}
	fs.storagePath = "./data_test/"

	// Create the test folder if it does not exist
	if err := os.MkdirAll(fs.storagePath, os.ModePerm); err != nil {
		t.Fatalf("Error creating test folder: %v", err)
	}

	// Cleanup the test folder after the test
	t.Cleanup(func() {
		files, err := os.ReadDir(fs.storagePath)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			os.Remove(fs.storagePath + file.Name())
		}
	})

	return fs
}

func TestSaveService(t *testing.T) {
	fs := setupTest(t)

	service := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "123",
			Name: "TestService",
		},
	}

	// Save the service
	err := fs.SaveService(&service)
	if err != nil {
		t.Fatalf("Error saving service: %v", err)
	}

	// Verify that the file was created
	if _, err := os.Stat(fs.storagePath + service.Info.Id + ".json"); os.IsNotExist(err) {
		t.Fatalf("Service file %s not found", service.Info.Id)
	}
}

func TestGetService(t *testing.T) {
	fs := setupTest(t)

	service := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "123",
			Name: "TestService",
		},
	}

	// Save the service
	if err := fs.SaveService(&service); err != nil {
		t.Fatalf("Error saving service: %v", err)
	}

	// Retrieve the service
	retrievedService, err := fs.GetService(service.Info.Id)
	if err != nil {
		t.Fatalf("Error retrieving service: %v", err)
	}

	// Verify that the retrieved service is correct
	if retrievedService.Info.Id != service.Info.Id {
		t.Errorf("Retrieved service ID does not match. Expected %s, got %s", service.Info.Id, retrievedService.Info.Id)
	}
}

func TestGetActiveServices(t *testing.T) {
	fs := setupTest(t)

	// Save two sample services
	service1 := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "123",
			Name: "TestService1",
		},
	}
	service2 := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "456",
			Name: "TestService2",
		},
	}

	if err := fs.SaveService(&service1); err != nil {
		t.Fatalf("Error saving service 1: %v", err)
	}
	if err := fs.SaveService(&service2); err != nil {
		t.Fatalf("Error saving service 2: %v", err)
	}

	// Retrieve active services
	services, err := fs.GetActiveServices()
	if err != nil {
		t.Fatalf("Error retrieving active services: %v", err)
	}

	// Verify that both services are retrieved
	if len(services) != 2 {
		t.Fatalf("Expected 2 active services, but found %d", len(services))
	}

	// Verify that the retrieved services match the saved ones
	var foundService1, foundService2 bool
	for _, srv := range services {
		if srv.Info.Id == service1.Info.Id {
			foundService1 = true
		}
		if srv.Info.Id == service2.Info.Id {
			foundService2 = true
		}
	}

	if !foundService1 {
		t.Errorf("Service %s was not found among active services", service1.Info.Id)
	}
	if !foundService2 {
		t.Errorf("Service %s was not found among active services", service2.Info.Id)
	}

	// Cleanup: delete services
	if err := fs.DeleteService(service1.Info.Id); err != nil {
		t.Fatalf("Error deleting service 1: %v", err)
	}
	if err := fs.DeleteService(service2.Info.Id); err != nil {
		t.Fatalf("Error deleting service 2: %v", err)
	}
}

func TestDeleteService(t *testing.T) {
	fs := setupTest(t)

	service := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "123",
			Name: "TestService",
		},
	}

	// Save the service
	if err := fs.SaveService(&service); err != nil {
		t.Fatalf("Error saving service: %v", err)
	}

	// Delete the service
	if err := fs.DeleteService(service.Info.Id); err != nil {
		t.Fatalf("Error deleting service: %v", err)
	}

	// Verify that the service file is deleted
	if _, err := os.Stat(fs.storagePath + service.Info.Id + ".json"); !os.IsNotExist(err) {
		t.Fatalf("Service file %s was not deleted", service.Info.Id)
	}
}

func TestIsService(t *testing.T) {
	fs := setupTest(t)

	service := models.PublicService{
		Info: models.ServiceInfo{
			Id:   "123",
			Name: "TestService",
		},
	}

	// Save the service
	if err := fs.SaveService(&service); err != nil {
		t.Fatalf("Error saving service: %v", err)
	}

	// Verify that the service exists
	if exists := fs.IsService(service.Info.Id); !exists {
		t.Fatalf("Service %s should exist", service.Info.Id)
	}

	// Delete the service
	if err := fs.DeleteService(service.Info.Id); err != nil {
		t.Fatalf("Error deleting service: %v", err)
	}

	// Verify that the service no longer exists
	if exists := fs.IsService(service.Info.Id); exists {
		t.Fatalf("Service %s should not exist after deletion", service.Info.Id)
	}
}

