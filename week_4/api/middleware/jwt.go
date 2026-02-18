package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. extract token from Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		// 2. parse "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
		}

		tokenString := parts[1]

		// 3. parse and validate the token
		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // this verifies signature + expiry in one call
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // verfifies the signing method
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		// 4. extract claims and store in context (but this is optional)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user", claims["sub"])
		}

		// 5. continue to the next handler
		return next(c)
	}
}
