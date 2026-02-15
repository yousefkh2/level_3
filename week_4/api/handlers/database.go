package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"time"
)

// for a func tro be used as a handler in Echo, it must satisify the signature (echo.Context) and return an error
func CreateDatabase(c echo.Context) error {

	// 1. get the service from context
	appCtx := c.Get("app").(*app.App)

	// 2. parse request body
	var req models.CreateDatabaseRequest
	/* Echo docs:
	Parsing request data in Echo is done with a process called binding.
	There are different ways to bind.
	https://echo.labstack.com/docs/binding
	*/
	// you give Bind the destination (the struct), even though it's working on the JSON
	// So, where is the source?
	// The input is hidden inside the c (the echo.Context).
	// This is also why you have to pass a pointer (&u) rather than just u.
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// 3. call the service
	cluster, err := appCtx.DBService.CreateDatabase(c.Request().Context(), req.Name, req.Instances, req.Storage)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "database already exists"})
		}
		if errors.IsInvalid(err) { // bad spec, like negative instances
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid database configuration"})
		}
		if errors.IsUnauthorized(err) { //RBAC issues
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized to create database"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 4. build response
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

	// call the service
	clusters, err := appCtx.DBService.ListDatabases(c.Request().Context())
	if err != nil {
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
	name := c.Param("name") // from URL like /database/:name

	cluster, err := appCtx.DBService.GetDatabase(c.Request().Context(), name)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	connectionInfo := models.ConnectionInfo{
		Host:     fmt.Sprintf("%s-rw.default.svc.cluster.local", cluster.Name), //
		Port:     5432,                                                         //
		Database: "app",                                                        // CNPG's default database name
		Username: "app",                                                        // CNPG's default username
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
			return c.NoContent(http.StatusNoContent) // idempotent (already gone)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func UpdateDatabase(c echo.Context) error {
	appCtx := c.Get("app").(*app.App)
	name := c.Param("name")

	var req models.UpdateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// at least one field must be provided
	if req.Instances == nil && req.Storage == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no fields to update"})
	}

	// call the service
	cluster, err := appCtx.DBService.UpdateDatabase(c.Request().Context(), name, req.Instances, req.Storage)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
		}
		if errors.IsInvalid(err) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
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
