package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
)

var ErrParameterMissing = errors.New("parameter is missing")

type Context struct {
	echo.Context
	middleware.SessionContext
}
