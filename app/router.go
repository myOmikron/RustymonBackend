package app

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/handler"
	"gorm.io/gorm"
)

func defineRoutes(e *echo.Echo, config *configs.RustymonConfig, db *gorm.DB) {
	// Either one, GET or POST are fine
	serverinfo := handler.ServerInfoHandler{
		DB:     db,
		Config: config,
	}
	e.GET("/serverinfo", serverinfo.Serverinfo())

	account := handler.AccountHandler{
		DB:     db,
		Config: config,
	}
	e.GET("/logout", account.Logout())
	e.POST("/logout", account.Logout())

	e.POST("/login", account.Login())
	e.POST("/register", account.Register())
	e.POST("/resetPassword", account.ResetPassword())
}
