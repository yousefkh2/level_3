package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
)

// for a func tro be used as a handler in Echo, it must satisify the signature (echo.Context) and return an error
func CreateDatabase(c echo.Context) error {

	// 1. get the service from context
	appCtx := c.Get("app").(*app.App)

	// 2. parse request body
	var req models.CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// 3. call the service
	cluster, err := appCtx.DBService.CreateDatabase(c.Request().Context(), req.Name, req.Instances, req.Storage)
	if err != nil {
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
