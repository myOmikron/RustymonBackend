package handler

import (
	"github.com/myOmikron/echotools/db"
	"github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
	"net/http"
)

func GetUsers(c *Context) error {
	var u []utilitymodels.User
	db.DB.Find(&u)
	return c.JSON(http.StatusOK, &utility.JsonResponse{
		Success: true,
		Data:    u,
	})
}
