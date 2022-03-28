package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/echotools/logging"
	"github.com/myOmikron/echotools/middleware"
	"gorm.io/gorm"
	"net/http"
	"runtime"
	"time"
)

func InitializeMiddleware(e *echo.Echo, db *gorm.DB) {
	log := logging.GetLogger("middleware")
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				log.Error(err.Error())
			}
			log.Infof(
				"%d %s %s %v - %s %s",
				c.Response().Status, c.Request().Method, c.Request().RequestURI, time.Now().Sub(start),
				c.RealIP(), c.Request().Header.Get("User-Agent"),
			)
			return nil
		}
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					if r == http.ErrAbortHandler {
						panic(r)
					}
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					stack := make([]byte, 4<<10) // 4 KB Stack size
					var length int
					if logging.GetLogLevel() == logging.DEBUG {
						length = runtime.Stack(stack, true)
					} else {
						length = runtime.Stack(stack, false)
					}
					stack = stack[:length]

					log.Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
				}
			}()
			return next(c)
		}
	})
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
