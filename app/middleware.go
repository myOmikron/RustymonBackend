package app

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/logging"
	"github.com/myOmikron/echotools/middleware"
	"gorm.io/gorm"
	"net/url"
	"time"
)

func InitializeMiddleware(e *echo.Echo, db *gorm.DB, config *configs.RustymonConfig) {
	log := logging.GetLogger("middleware")
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(middleware.Panic(log))
	e.Use(middleware.Logging(log))

	// Security unpacking
	allowedHosts := &middleware.SecurityConfig{
		AllowedHosts:            []middleware.AllowedHost{},
		UseForwardedProtoHeader: config.Server.UseForwardedProtoHeader,
	}
	for _, allowedHost := range config.Server.AllowedHosts {
		u, _ := url.Parse(allowedHost)
		https := false
		if u.Scheme == "https" {
			https = true
		}
		allowedHosts.AllowedHosts = append(allowedHosts.AllowedHosts, middleware.AllowedHost{
			Host:  u.Host,
			Https: https,
		})
	}
	e.Use(middleware.Security(log, allowedHosts))
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
