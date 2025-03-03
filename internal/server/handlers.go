
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingHandler responds with a simple "pong" message
func PingHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}


