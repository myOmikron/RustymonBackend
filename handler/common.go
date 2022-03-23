package handler

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
)

var json = jsoniter.Config{
	EscapeHTML:    true,
	CaseSensitive: true,
}.Froze()

var ErrParameterMissing = errors.New("parameter is missing")

type Context struct {
	echo.Context
	middleware.SessionContext
}
