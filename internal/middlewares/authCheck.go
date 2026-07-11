package middlewares

import (
	"gotickets/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// extract token from authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
					"message": "Missing authorization header",
				})
			}
			//check bearer schema
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
					"message": "Invalid authorization header format",
				})
			}
			tokenString := parts[1]
			// validation token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
					"message": "Invalid or expired token",
				})
			}
			// store user info in context from handler
			c.Set("user_id", claims.UserId)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)

			return next(c)
		}
	}
}
