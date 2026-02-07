package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yousefkh2/level_3/week_4/api/config"
	"k8s.io/client-go/kubernetes"
	"log"
)

type App struct {
	K8sClient *kubernetes.Clientset
	// Later: database connections, config, etc.
}

func main() {
	k8sClient, err := config.NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		K8sClient: k8sClient,
	}

	_ = app
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes will go here

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
