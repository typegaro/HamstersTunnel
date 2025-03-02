package service

import (
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
)

func GeneratePublicService(bannedPort []string) (error,models.PublicService) {
    //TODO: Generate and check port  
	return nil,models.PublicService{
		PublicHTTP: 8032,
		PublicTCP:  2023,
		PublicUDP:  4021,
	}
}
