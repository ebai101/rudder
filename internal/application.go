package internal

import (
	"rudder/internal/clients"
	"rudder/internal/config"
	"rudder/internal/database"
	"rudder/internal/handlers"
	"rudder/internal/repositories"
	"rudder/internal/services"
	"rudder/util/routing"
	"rudder/util/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	E        *echo.Echo
	DB       *database.DBConnection
	Config   *config.AppConfig
	Args     config.Args
	SrvFin   *services.FinancialService
	SrvSFIN  *services.SimpleFINService
	SrvAc    *services.AutocatService
	SrvSched *services.SchedService
	HTxn     *handlers.TransactionsHandlers
	HAcc     *handlers.AccountsHandlers
	HAcat    *handlers.AutocatHandlers
}

func NewApplication(c *config.AppConfig, args config.Args) (*Application, error) {
	e := bootstrapEcho()

	db, err := database.NewDBConnection(c)
	if err != nil {
		return nil, err
	}

	clientSfin := clients.NewSimpleFINClient(c)
	repoSfin := repositories.NewSimpleFINRepository(db)
	srvSfin := services.NewSimpleFINService(c, clientSfin, repoSfin)

	repoAcat := repositories.NewAutocatRepository(db)
	srvAcat := services.NewAutocatService(repoAcat, srvSfin)
	hAcat := handlers.NewAutocatHandlers(srvAcat)

	repoFin := repositories.NewFinancialRepository(db)
	srvFin := services.NewFinancialService(repoFin)
	hTxn := handlers.NewTransactionsHandlers(srvFin)
	hAcc := handlers.NewAccountsHandlers(srvFin)

	sched := services.NewSchedService(c, args, srvSfin)

	app := &Application{
		E:        e,
		DB:       db,
		Config:   c,
		Args:     args,
		SrvFin:   srvFin,
		SrvSFIN:  srvSfin,
		SrvAc:    srvAcat,
		SrvSched: sched,
		HTxn:     hTxn,
		HAcc:     hAcc,
		HAcat:    hAcat,
	}

	return app, nil
}

func bootstrapEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	routing.SetupRouter(e)
	template.NewTemplateRenderer(e)

	e.Static("/static", "assets")

	return e
}
