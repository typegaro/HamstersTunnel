package server

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggerMiddleware logs HTTP request information including method, path and duration
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c) // Execute next handler
		duration := time.Since(start)

		log.Printf("%s %s | %v", c.Request().Method, c.Request().URL.Path, duration)
		return err
	}
}
