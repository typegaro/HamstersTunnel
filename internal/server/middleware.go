package server

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggerMiddleware registra informazioni sulle richieste HTTP
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c) // Esegue il prossimo handler
		duration := time.Since(start)

		log.Printf("%s %s | %v", c.Request().Method, c.Request().URL.Path, duration)
		return err
	}
}
