package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
	"go.uber.org/zap"
)

func Login(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)

	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		appCtx.Logger.Warn("failed to parse login request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	expectedUser := os.Getenv("API_USERNAME")
	expectedPass := os.Getenv("API_PASSWORD")

	if req.Username != expectedUser || req.Password != expectedPass {
		appCtx.Logger.Warn("login failed: invalid credentials", zap.String("username", req.Username))
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		appCtx.Logger.Error("failed to sign JWT token", zap.String("username", req.Username), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to sign token"})
	}

	appCtx.Logger.Info("user logged in successfully", zap.String("username", req.Username))

	return c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString})
}
