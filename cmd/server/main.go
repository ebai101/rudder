package main

import (
	"log"
	"rudder/internal/config"
	"rudder/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// args := config.ParseArgs()

	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(appConfig)

	// if err := config.LoadConfig(appConfig); err != nil {
	// 	log.Fatal(err)
	// }

	app := echo.New()

	app.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	app.Static("/", "assets")
	app.Use(middleware.Logger())

	// client := clients.NewSimpleFINClient(appConfig)

	// db, err := database.NewDBConnection(appConfig)
	// if err != nil {
	// 	app.Logger.Fatalf("failed to create db connection: %s", err)
	// }

	// repo := repositories.NewSimpleFINRepository(db)
	// s := services.NewSimpleFINService(appConfig, client, repo)
	// fmt.Println(s)

	app.Logger.Fatal(app.Start(":4040"))
}
