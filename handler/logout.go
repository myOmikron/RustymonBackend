package handler

import (
	"github.com/myOmikron/echotools/middleware"
	u "github.com/myOmikron/echotools/utility"
)

func Logout(c *Context) error {
	if err := middleware.Logout(c); err != nil {
		return c.JSON(500, u.JsonResponse{Error: err.Error()})
	}
	return c.JSON(200, u.JsonResponse{Success: true})
}
