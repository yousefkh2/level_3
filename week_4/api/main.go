package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/config"
	"github.com/yousefkh2/level_3/week_4/api/handlers"
	mw "github.com/yousefkh2/level_3/week_4/api/middleware"
	"github.com/yousefkh2/level_3/week_4/api/services"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	k8sClient, err := config.NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}

	application := &app.App{
		K8sClient: k8sClient,
		DBService: services.NewDatabaseService(k8sClient),
		Logger:    logger,
		Metrics:   app.NewMetrics(),
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.PATCH, echo.DELETE, echo.OPTIONS},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", application)
			return next(c)
		}
	})

	e.Use(mw.Observability(application))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	e.POST("/auth/login", handlers.Login)

	e.POST("/databases", handlers.CreateDatabase, mw.JWTAuth)
	e.GET("/databases", handlers.ListDatabases, mw.JWTAuth)

	e.GET("/databases/:name", handlers.GetDatabase, mw.JWTAuth)
	e.PATCH("/databases/:name", handlers.UpdateDatabase, mw.JWTAuth)
	e.DELETE("/databases/:name", handlers.DeleteDatabase, mw.JWTAuth)
	logger.Info("starting API server", zap.String("addr", ":8080"))
	e.Logger.Fatal(e.Start(":8080"))
}
