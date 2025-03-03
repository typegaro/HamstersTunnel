package service

import (
	"net/http"
    
	"github.com/labstack/echo/v4"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
)

// Handler for creating a new service
func (s *ServiceManager) HandlerNewService(c echo.Context) error {

	save := c.QueryParam("save")
    if save != "true" && save != "false" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid save value"})
	}

	var req models.NewServiceReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

    //TODO: Handle the error 
	publicService,_ := GeneratePublicService(req)

	// Salva il servizio su file
	//err := SaveServiceToFile(service)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save service"})
	//}

	// Restituisci il PublicService come risposta
	return c.JSON(http.StatusOK, map[string]string{"service_id": publicService.Info.Id})
}
