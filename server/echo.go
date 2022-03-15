package server

import (
	"RustymonBackend/handler"
	"RustymonBackend/models"
	"RustymonBackend/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var c = utils.ChangeContext

func StartServer() {
	// Echo instance
	e := echo.New()

	// Set debug level
	e.Logger.SetLevel(log.DEBUG)

	// Open DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// Migrate
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(err.Error())
	}

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &utils.Context{Context: c, DB: db}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/info", c(handler.Info))
	e.GET("/users", c(handler.GetUsers))
	e.POST("/register", c(handler.Register))

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
