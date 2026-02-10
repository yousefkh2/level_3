package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/config"
	"github.com/yousefkh2/level_3/week_4/api/handlers"
	"github.com/yousefkh2/level_3/week_4/api/services"
)

func main() {
	k8sClient, err := config.NewKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{
		K8sClient: k8sClient,
		DBService: services.NewDatabaseService(k8sClient),
	}

	_ = app
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// store app in Echo's context for use in handlers
	/* The "Why" Behind the Pattern
	By default, an Echo handler looks like this: func(c echo.Context) error. It only knows about the HTTP request. It has no idea your K8sClient exists.
	By using c.Set("app", app), you are attaching your backend logic to the Context of that specific request. It’s like giving every waiter (handler) a master key (the App struct) to the kitchen (the Database/K8s).
	*/

	/*
		The reason that last one is harder to remember is that middleware.Logger() is a pre-built function that returns the middleware for you.
		In this one, we are writing the middleware ourselves.
	*/
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc { // e.Use(...) tells Echo: "Run this for every request."
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	// routes will go here

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	e.POST("/databases", handlers.CreateDatabase)
	e.GET("/databases", handlers.ListDatabases)

	e.GET("/databases/:name", handlers.GetDatabase)
	e.DELETE("/databases/:name", handlers.DeleteDatabase)
	e.Logger.Fatal(e.Start(":8080"))
}
