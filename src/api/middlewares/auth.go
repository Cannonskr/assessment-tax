package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AdminAuthMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "adminTax" && password == "admin!" {
			return true, nil
		}
		return false, nil
	})
}
