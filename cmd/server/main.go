package main 

import (
	"github.com/labstack/echo/v4"
	"github.com/typegaro/HamstersTunnel/internal/server"
)

func main() {
	e := echo.New()
    s := server.NewServer()

	// Configure routes
	server.SetupRoutes(e, s)
    e.Use(server.LoggerMiddleware)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
