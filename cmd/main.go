package main

import (
	"log"
	"rudder/config"
	"rudder/proc"
	"rudder/resource"
)

func main() {
	log.Println("Parsing arguments...")
	args := config.ParseArgs()

	log.Println("Loading config...")
	var appConfig config.AppConfig
	if err := config.LoadConfig(&appConfig); err != nil {
		log.Fatal(err)
	}

	log.Println("Setting up database...")
	db, err := resource.OpenDatabase(&appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Setting up API...")
	sfinAPI := resource.SimpleFINAPI{
		Config: &appConfig,
	}

	log.Println("Done setting up.")

	if args.Scheduled {
		log.Println("Starting scheduled tasks")
		err := proc.StartScheduler(&appConfig, &db, &sfinAPI, &args)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Doing a one-time update")
		if err := proc.Update(&appConfig, &db, &sfinAPI, &args, args.DaysToFetch); err != nil {
			log.Fatal(err)
		}
	}
}
