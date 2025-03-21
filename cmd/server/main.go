package main 

import (
	"github.com/labstack/echo/v4"
	"github.com/typegaro/HamstersTunnel/internal/server"
)

func main() {
	e := echo.New()
    s := server.NewServer()
    s.Init()

	server.SetupRoutes(e, s)
    e.Use(server.LoggerMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
