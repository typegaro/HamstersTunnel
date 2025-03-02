package server

import "github.com/labstack/echo/v4"


func SetupRoutes(e *echo.Echo) {
	e.GET("/ping", PingHandler)
	e.POST("/service/:service", NewServiceHandler) 
}
