package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
	"github.com/yousefkh2/level_3/week_4/api/services"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultAuditLookbackHours = 24

func parseAuditLookback(c echo.Context) (time.Duration, error) {
	hoursParam := c.QueryParam("hours")
	if hoursParam == "" {
		return defaultAuditLookbackHours * time.Hour, nil
	}

	hours, err := strconv.Atoi(hoursParam)
	if err != nil || hours <= 0 {
		return 0, fmt.Errorf("invalid hours query param: must be a positive integer")
	}

	return time.Duration(hours) * time.Hour, nil
}

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

	appCtx.LogAuditEvent(c.Request().Context(), app.AuditActionCreate, req.Name,
		zap.Int("instances", req.Instances),
		zap.String("storage", req.Storage),
	)

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

	password := ""
	var appSecret corev1.Secret
	secretName := fmt.Sprintf("%s-app", cluster.Name)
	secretErr := appCtx.K8sClient.Get(c.Request().Context(), client.ObjectKey{
		Name:      secretName,
		Namespace: cluster.Namespace,
	}, &appSecret)
	if secretErr != nil {
		appCtx.Logger.Warn("failed to get database credentials secret", zap.String("name", name), zap.String("secret", secretName), zap.Error(secretErr))
	} else {
		password = string(appSecret.Data["password"])
	}

	connectionInfo := models.ConnectionInfo{
		Host:     fmt.Sprintf("%s-rw.%s.svc.cluster.local", cluster.Name, cluster.Namespace),
		Port:     5432,
		Database: "app",
		Username: "app",
		Password: password,
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

	appCtx.LogAuditEvent(c.Request().Context(), app.AuditActionDelete, name)

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

	fields := []zap.Field{}
	if req.Instances != nil {
		fields = append(fields, zap.Int("instances", *req.Instances))
	}
	if req.Storage != nil {
		fields = append(fields, zap.String("storage", *req.Storage))
	}
	appCtx.LogAuditEvent(c.Request().Context(), app.AuditActionUpdate, name, fields...)

	return c.JSON(http.StatusOK, response)
}

func GetDatabaseLogs(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)
	name := c.Param("name")

	since, err := parseAuditLookback(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	events, err := appCtx.LokiService.GetAuditLogs(c.Request().Context(), name, since)
	if err != nil {
		appCtx.Logger.Error("failed to fetch audit logs", zap.String("database", name), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch audit logs"})
	}

	// Return empty array instead of null for consistency
	if events == nil {
		events = make([]services.AuditEvent, 0)
	}

	appCtx.Logger.Info("audit logs retrieved", zap.String("database", name), zap.Int("count", len(events)))

	return c.JSON(http.StatusOK, events)
}

func GetGlobalAuditLogs(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)

	since, err := parseAuditLookback(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	events, err := appCtx.LokiService.GetGlobalAuditLogs(c.Request().Context(), since)
	if err != nil {
		appCtx.Logger.Error("failed to fetch global audit logs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch audit logs"})
	}

	if events == nil {
		events = make([]services.AuditEvent, 0)
	}

	appCtx.Logger.Info("global audit logs retrieved", zap.Int("count", len(events)))

	return c.JSON(http.StatusOK, events)
}
