package main

import (
	"context"
	"log"
	"os"
	"rudder/internal"
	"rudder/internal/config"
	"rudder/internal/handlers"
)

func update(app *internal.Application, args config.Args) error {
	err := app.SrvSFIN.SyncSimpleFIN(
		context.Background(),
		args.UseCached,
		args.SaveCached,
		args.DaysToFetch,
	)
	if err != nil {
		return err
	}

	_, err = app.SrvCat.CategorizeTransactions(context.Background())
	if err != nil {
		return err
	}

	return nil
}

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

	if args.Update {
		err := update(app, args)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	handlers.RegisterRoutes(app.E, app.HIns, app.HTxn, app.HAcc, app.HCat)
	app.SrvSched.Start()
	app.E.Logger.Fatal(app.E.Start(":4040"))
}
