package utils

import (
	"RustymonBackend/middleware/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Context struct {
	echo.Context
	DB      *gorm.DB
	Session *session.Session
}

func ChangeContext(f func(c Context) error) echo.HandlerFunc {
	return func(context echo.Context) error {
		return f(*context.(*Context))
	}
}
