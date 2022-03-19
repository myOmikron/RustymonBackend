package handler

import (
	"encoding/json"
	"errors"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/middleware"
	u "github.com/myOmikron/echotools/utility"
	"io/ioutil"
)

var ErrLoginFailed = errors.New("login failed")

type LoginForm struct {
	Username string
	Password string
}

func Login(c *Context) error {
	var l LoginForm
	b, _ := ioutil.ReadAll(c.Request().Body)

	if err := json.Unmarshal(b, &l); err != nil {
		return c.JSON(400, u.JsonResponse{Error: err.Error()})
	}

	if l.Username == "" || l.Password == "" {
		return c.JSON(400, u.JsonResponse{
			Error: ErrParameterMissing.Error(),
		})
	}

	user, err := auth.Authenticate(l.Username, l.Password)
	if err != nil || user == nil {
		c.Logger().Info(err)
		return c.JSON(403, u.JsonResponse{Error: ErrLoginFailed.Error()})
	} else {
		if err := middleware.Login(user, c); err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}
	}

	return c.JSON(200, u.JsonResponse{Success: true})
}
