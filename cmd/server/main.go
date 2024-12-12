package main

import (
	"log"
	"rudder/internal"
	"rudder/internal/config"
	"rudder/internal/handlers"
)

func main() {
	args := config.ParseArgs()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := internal.NewApplication(c, args)
	if err != nil {
		log.Fatal(err)
	}

	handlers.RegisterRoutes(app)

	app.Sched.Start()

	app.E.Logger.Fatal(app.E.Start(":4040"))
}
