package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/models"
)

func Login(c echo.Context) error {
	// 1. parse request body
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// 2. validate credentials against env vars
	expectedUser := os.Getenv("API_USERNAME")
	expectedPass := os.Getenv("API_PASSWORD")

	if req.Username != expectedUser || req.Password != expectedPass {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	// 3. sign the JWT
	jwtSecret := os.Getenv("JWT_SECRET")

	// Header: SigningMethodHS256 (algorithm)
	// Payload: MapClaims (the data)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,                          // who the token is for
		"exp": time.Now().Add(24 * time.Hour).Unix(), // expiry
		"iat": time.Now().Unix(),                     // issued at
	})

	// Signature: HMAC of header+payload using the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to sign token"})
	}

	// 4. return the token
	return c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString})
}
