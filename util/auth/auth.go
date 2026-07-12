package auth

import "github.com/labstack/echo/v4"

func UserID(c echo.Context) int {
	id, _ := c.Get("id").(int)
	return id
}
