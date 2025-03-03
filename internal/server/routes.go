package server

import (
    "github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo,s *Server) {
	e.GET("/ping", PingHandler)
	e.POST("/service", s.ServiceManager.HandlerNewService) 
}
