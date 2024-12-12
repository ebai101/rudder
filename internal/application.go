package internal

import (
	"rudder/internal/clients"
	"rudder/internal/config"
	"rudder/internal/database"
	"rudder/internal/repositories"
	"rudder/internal/services"
	"rudder/util/routing"
	"rudder/util/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	E      *echo.Echo
	DB     *database.DBConnection
	Config *config.AppConfig
	Args   config.Args
	SFIN   *services.SimpleFINService
	AC     *services.AutocatService
	Sched  *services.SchedService
}

func NewApplication(c *config.AppConfig, args config.Args) (*Application, error) {
	e := bootstrapEcho()

	db, err := database.NewDBConnection(c)
	if err != nil {
		return nil, err
	}

	sfinC := clients.NewSimpleFINClient(c)
	sfinR := repositories.NewSimpleFINRepository(db)
	sfin := services.NewSimpleFINService(c, sfinC, sfinR)

	acR := repositories.NewAutocatRepository(db)
	ac := services.NewAutocatService(acR, sfin)

	sched := services.NewSchedService(c, args, sfin)

	app := &Application{
		E:      e,
		DB:     db,
		Config: c,
		Args:   args,
		SFIN:   sfin,
		AC:     ac,
		Sched:  sched,
	}

	return app, nil
}

func bootstrapEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	routing.SetupRouter(e)
	template.NewTemplateRenderer(e)

	e.Static("/", "assets")

	return e
}
