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
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid save value"})
	}

	var req models.NewServiceReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	publicService, err := GeneratePublicService(req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "failed to generate service"},
		)
	}

	if save {
		if err := s.memory.AddService(&publicService); err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "failed to save service"},
			)
		}
	}
	serviceResponse := models.NewServiceRes{
		Id: publicService.Id,
	}

	if publicService.HTTP != nil {
		serviceResponse.HTTP = ":" + publicService.HTTP.Client
	}
	if publicService.TCP != nil {
		serviceResponse.TCP = ":" + publicService.TCP.Client
	}
	if publicService.UDP != nil {
		serviceResponse.UDP = ":" + publicService.UDP.Client
	}

	return c.JSON(http.StatusOK, serviceResponse)
}

func (s *ServiceManager) HandlerStartService(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing service ID"})
	}
	srv := s.memory.GetService(id)
	srv.Active = true
	s.memory.EditService(srv)

	return c.JSON(http.StatusOK, map[string]string{"message": "Service started", "id": id})
}

func (s *ServiceManager) HandlerStopService(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing service ID"})
	}
	srv := s.memory.GetService(id)
	srv.Active = true
	s.memory.EditService(srv)

	return c.JSON(http.StatusOK, map[string]string{"message": "Service stopped", "id": id})
}

func (s *ServiceManager) HandlerRemoveService(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing service ID"})
	}
	s.memory.RemoveService(id)

	return c.JSON(http.StatusOK, map[string]string{"message": "Service removed", "id": id})
}
