package main

import (
	"log"
	"rudder/internal"
	"rudder/internal/config"
	"rudder/internal/database"
	"rudder/internal/handlers"
	"rudder/util/routing"
	"rudder/util/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	args := config.ParseArgs()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := bootstrap()
	db, err := database.NewDBConnection(c)
	if err != nil {
		log.Fatal(err)
	}

	app := &internal.Application{
		E:      e,
		DB:     db,
		Config: c,
		Args:   args,
	}

	handlers.RegisterRoutes(app)

	e.Logger.Fatal(e.Start(":4040"))
}

func bootstrap() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	routing.SetupRouter(e)
	template.NewTemplateRenderer(e)

	e.Static("/", "assets")

	return e
}
