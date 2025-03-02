
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
    "github.com/typegaro/HamstersTunnel/internal/service"
    
)

// PingHandler risponde con "pong"
func PingHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

// Handler for creating a new service
func NewServiceHandler(c echo.Context) error {
	name := c.Param("service")

	save := c.QueryParam("save")
	if save != "true" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing save=true"})
	}

	var req models.LocalService
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

    //TODO: Add the list of banned port
	_,publicService := service.GeneratePublicService([]string{})

	// Crea il servizio completo
	service := &models.Service{
		Name:     name,
		LService: req,
		PService: publicService,
	}

	// Salva il servizio su file
	//err := SaveServiceToFile(service)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save service"})
	//}

	// Restituisci il PublicService come risposta
	return c.JSON(http.StatusOK, service.PService)
}
