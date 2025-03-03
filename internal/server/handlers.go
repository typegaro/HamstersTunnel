
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingHandler risponde con "pong"
func PingHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}


