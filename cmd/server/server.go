package main

import (
	"github.com/labstack/echo/v4"
	"github.com/typegaro/HamstersTunnel/internal/server"
)

func main() {
	e := echo.New()

	// Configura le rotte
	server.SetupRoutes(e)
    e.Use(server.LoggerMiddleware)

	// Avvia il server
	e.Logger.Fatal(e.Start(":8080"))
}
