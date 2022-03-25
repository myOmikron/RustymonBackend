package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
)

type Context struct {
	echo.Context
	middleware.SessionContext
}
