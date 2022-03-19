package server

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/db"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/driver/sqlite"
	"time"
)

func StartServer() {
	// Echo instance
	e := echo.New()

	// Set debug level
	e.Logger.SetLevel(log.DEBUG)

	// Initialize DB
	db.Initialize(
		sqlite.Open("test.db"),
		&utilitymodels.Session{},
	)

	// Set session middleware config

	// Middleware
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Gzip())
	f := false
	age := time.Hour * 24
	e.Use(middleware.Session(&middleware.SessionConfig{
		Secure:         &f,
		CookieAge:      &age,
		DisableLogging: true,
	}))

	// Routes
	e.GET("/info", middleware.Wrap(handler.Info))
	e.POST("/info", middleware.Wrap(handler.Info))

	e.POST("/register", middleware.Wrap(handler.Register))
	e.POST("/login", middleware.Wrap(handler.Login))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
