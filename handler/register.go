package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/echotools/db"
	u "github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
	"io/ioutil"
	"net/mail"
)

var ErrUsernameOrEmailTaken = errors.New("username or email already exists")

type RegisterForm struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	TrainerName string `json:"trainer_name"`
}

func Register(c *Context) error {
	var f RegisterForm
	b, _ := ioutil.ReadAll(c.Request().Body)

	if err := json.Unmarshal(b, &f); err != nil {
		return c.JSON(400, u.JsonResponse{Error: err.Error()})
	}
	if f.Username == "" || f.Password == "" || f.Email == "" || f.TrainerName == "" {
		return c.JSON(400, u.JsonResponse{Error: ErrParameterMissing.Error()})
	}

	if address, err := mail.ParseAddress(f.Email); err != nil {
		return c.JSON(400, u.JsonResponse{Error: fmt.Sprintf("No valid mail provided: %s", err.Error())})
	} else {
		f.Email = address.Address
	}

	var count int64
	db.DB.Find(&utilitymodels.User{}, "username = ? OR email = ?", f.Username, f.Email).Count(&count)
	if count > 0 {
		return c.JSON(409, u.JsonResponse{Error: ErrUsernameOrEmailTaken.Error()})
	}

	user, err := db.CreateUser(f.Username, f.Password, f.Email, true)
	if err != nil {
		return c.JSON(500, u.JsonResponse{Error: err.Error()})
	}
	player := models.Player{
		User:        *user,
		TrainerName: f.TrainerName,
	}
	if err := db.DB.Create(&player).Error; err != nil {
		return c.JSON(500, u.JsonResponse{Error: err.Error()})
	}

	return c.JSON(200, u.JsonResponse{Success: true})
}
