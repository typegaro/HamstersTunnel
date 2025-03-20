package server

import (
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all the API endpoints
func SetupRoutes(e *echo.Echo, s *Server) {
	e.GET("/ping", PingHandler)
	e.POST("/service", s.ServiceManager.HandlerNewService)

	e.POST("/service/:id/start", s.ServiceManager.HandlerStartService)
	e.PUT("/service/:id/stop", s.ServiceManager.HandlerStopService)
	e.DELETE("/service/:id", s.ServiceManager.HandlerRemoveService)
	//e.GET("/service/:id", s.ServiceManager.HandlerGetServiceStatus)
}
