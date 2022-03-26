package app

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

func defineRoutes(e *echo.Echo, config *configs.RustymonConfig, db *gorm.DB, wp worker.Pool) {
	// Either one, GET or POST are fine
	serverinfo := handler.ServerInfoHandler{
		DB:     db,
		Config: config,
	}

	account := handler.AccountHandler{
		DB:         db,
		Config:     config,
		WorkerPool: wp,
	}
	
	e.GET("/serverinfo", serverinfo.Serverinfo())

	e.GET("/logout", account.Logout())
	e.POST("/logout", account.Logout())

	e.POST("/login", account.Login())
	e.POST("/register", account.Register())
	e.POST("/resetPasswordByUsername", account.ResetPasswordUsername())
	e.POST("/resetPasswordByEmail", account.ResetPasswordEmail())
}
