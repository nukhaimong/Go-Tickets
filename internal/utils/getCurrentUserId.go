package utils

import "github.com/labstack/echo/v5"

func GetCurrentUserID(c *echo.Context) (uint, bool) {
	userId, ok := c.Get("user_id").(uint)
	return userId, ok
}
