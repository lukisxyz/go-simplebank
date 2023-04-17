package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if len(authHeader) == 0 {
			return c.NoContent(http.StatusUnauthorized)
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return c.NoContent(http.StatusUnauthorized)
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			return c.NoContent(http.StatusUnauthorized)
		}

		accessToken := fields[1]
		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("payload", payload)
		return next(c)
	}
}
