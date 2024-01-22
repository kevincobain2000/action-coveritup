package pkg

import (
	"github.com/labstack/echo/v4"
)

// HasAuthorizationHeader is a middleware to allow access to fail server
func HasAuthorizationHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") == "" {
				return echo.NewHTTPError(401, "Authorization header is required")
			}

			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}
