package app

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/logging"
	"github.com/myOmikron/echotools/middleware"
	"gorm.io/gorm"
	"time"
)

func InitializeMiddleware(e *echo.Echo, db *gorm.DB) {
	log := logging.GetLogger("middleware")
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(middleware.Logging(log))
	e.Use(middleware.Panic(log))
	e.Use(echoMiddleware.Gzip())
	f := false
	age := time.Hour * 24
	e.Use(middleware.Session(
		db,
		log,
		&middleware.SessionConfig{
			Secure:         &f,
			CookieAge:      &age,
			DisableLogging: true,
		},
	))
}
