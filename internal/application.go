package internal

import (
	"rudder/internal/config"
	"rudder/internal/database"

	"github.com/labstack/echo/v4"
)

type Application struct {
	E      *echo.Echo
	DB     *database.DBConnection
	Config *config.AppConfig
	Args   config.Args
}
