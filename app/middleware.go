package app

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/middleware"
	"gorm.io/gorm"
	"os"
	"time"
)

func InitializeMiddleware(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format:           "",
		CustomTimeFormat: time.RFC1123Z,
		Output:           os.Stdout,
	}))
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Gzip())
	f := false
	age := time.Hour * 24
	e.Use(middleware.Session(
		db,
		&middleware.SessionConfig{
			Secure:         &f,
			CookieAge:      &age,
			DisableLogging: true,
		},
	))
}
