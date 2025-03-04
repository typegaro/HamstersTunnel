package service

import (
	"net/http"
    "strconv"
    
	"github.com/labstack/echo/v4"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
)

// Handler for creating a new service
func (s *ServiceManager) HandlerNewService(c echo.Context) error {

    save, err := strconv.ParseBool(c.QueryParam("save"))
    if err!=nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid save value"})
	}

	var req models.NewServiceReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	publicService,err := GeneratePublicService(req)
    if err!= nil{
        return c.JSON(http.StatusInternalServerError,map[string]string{"error": "failed to generate service"})
    }

    if save {
        if err := s.memory.SaveService(&publicService); err != nil {
		    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save service"})
	    }
    }

	return c.JSON(http.StatusOK, map[string]string{"service_id": publicService.Info.Id})
}
