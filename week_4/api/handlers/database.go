package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"time"
)

func CreateDatabase(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)

	var req models.CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		appCtx.Logger.Warn("failed to parse create database request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	cluster, err := appCtx.DBService.CreateDatabase(c.Request().Context(), req.Name, req.Instances, req.Storage)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			appCtx.Logger.Warn("database already exists", zap.String("name", req.Name))
			return c.JSON(http.StatusConflict, map[string]string{"error": "database already exists"})
		}
		if errors.IsInvalid(err) {
			appCtx.Logger.Warn("invalid database configuration", zap.String("name", req.Name), zap.Int("instances", req.Instances), zap.String("storage", req.Storage), zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid database configuration"})
		}
		if errors.IsUnauthorized(err) {
			appCtx.Logger.Error("unauthorized to create database", zap.String("name", req.Name), zap.Error(err))
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized to create database"})
		}
		appCtx.Logger.Error("failed to create database", zap.String("name", req.Name), zap.Int("instances", req.Instances), zap.String("storage", req.Storage), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := models.DatabaseResponse{
		Name: cluster.Name,
		Spec: models.DatabaseSpec{
			Instances: cluster.Spec.Instances,
			Storage:   cluster.Spec.StorageConfiguration.Size,
		},
		Status:    "creating",
		CreatedAt: time.Now(),
	}

	return c.JSON(http.StatusCreated, response)

}

func ListDatabases(c echo.Context) error {

	appCtx := c.Get("app").(*app.App)

	clusters, err := appCtx.DBService.ListDatabases(c.Request().Context())
	if err != nil {
		appCtx.Logger.Error("failed to list databases", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]models.DatabaseResponse, 0, len(clusters))
	for _, cluster := range clusters {
		responses = append(responses, models.DatabaseResponse{
			Name: cluster.Name,
			Spec: models.DatabaseSpec{
				Instances: cluster.Spec.Instances,
				Storage:   cluster.Spec.StorageConfiguration.Size,
			},
			Status:    string(cluster.Status.Phase),
			CreatedAt: cluster.CreationTimestamp.Time,
		})
	}

	return c.JSON(http.StatusOK, responses)

}

func GetDatabase(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)
	name := c.Param("name")

	cluster, err := appCtx.DBService.GetDatabase(c.Request().Context(), name)
	if err != nil {
		if errors.IsNotFound(err) {
			appCtx.Logger.Warn("database not found", zap.String("name", name))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
		}
		appCtx.Logger.Error("failed to get database", zap.String("name", name), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	connectionInfo := models.ConnectionInfo{
		Host:     fmt.Sprintf("%s-rw.default.svc.cluster.local", cluster.Name),
		Port:     5432,
		Database: "app",
		Username: "app",
	}

	response := models.DatabaseDetailResponse{
		Name: cluster.Name,
		Spec: models.DatabaseSpec{
			Instances: cluster.Spec.Instances,
			Storage:   cluster.Spec.StorageConfiguration.Size,
		},
		Status:     string(cluster.Status.Phase),
		CreatedAt:  cluster.CreationTimestamp.Time,
		Connection: connectionInfo,
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteDatabase(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)
	name := c.Param("name")

	err := appCtx.DBService.DeleteDatabase(c.Request().Context(), name)
	if err != nil {
		if errors.IsNotFound(err) {
			appCtx.Logger.Info("delete database: not found (idempotent)", zap.String("name", name))
			return c.NoContent(http.StatusNoContent)
		}
		appCtx.Logger.Error("failed to delete database", zap.String("name", name), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func UpdateDatabase(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)
	name := c.Param("name")

	var req models.UpdateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		appCtx.Logger.Warn("failed to parse update database request", zap.String("name", name), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Instances == nil && req.Storage == nil {
		appCtx.Logger.Warn("update database: no fields provided", zap.String("name", name))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no fields to update"})
	}

	cluster, err := appCtx.DBService.UpdateDatabase(c.Request().Context(), name, req.Instances, req.Storage)
	if err != nil {
		if errors.IsNotFound(err) {
			appCtx.Logger.Warn("update database: not found", zap.String("name", name))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
		}
		if errors.IsInvalid(err) {
			appCtx.Logger.Warn("update database: invalid configuration", zap.String("name", name), zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		appCtx.Logger.Error("failed to update database", zap.String("name", name), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := models.DatabaseResponse{
		Name: cluster.Name,
		Spec: models.DatabaseSpec{
			Instances: cluster.Spec.Instances,
			Storage:   cluster.Spec.StorageConfiguration.Size,
		},
		Status:    string(cluster.Status.Phase),
		CreatedAt: cluster.CreationTimestamp.Time,
	}

	return c.JSON(http.StatusOK, response)
}
