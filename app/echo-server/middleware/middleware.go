package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware2(jwtSign string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			signature := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: missing token"))
			}
			if signature[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: invalid token format"))
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSign), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: "+err.Error()))
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: invalid signing method"))
			}

			expAt, err := claim.GetExpirationTime()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: token has no expiry"))
			}
			if time.Now().After(expAt.Time) {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: token expired"))
			}

			userIDFloat, _ := claim["id"].(float64)
			c.Set("id", int(userIDFloat))

			return next(c)
		}
	}
}
