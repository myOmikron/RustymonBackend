package handler

import (
	"RustymonBackend/models"
	"RustymonBackend/utils"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
)

type RegisterForm struct {
	Username string
	Password string
	Nick     string
}

func Register(c utils.Context) error {
	var f RegisterForm
	b, _ := ioutil.ReadAll(c.Request().Body)

	if err := json.Unmarshal(b, &f); err != nil {
		return c.JSON(400, JsonResponse{Error: err})
	}
	if f.Username == "" || f.Nick == "" || f.Password == "" {
		return c.JSON(400, JsonResponse{Message: "Parameter username, nick and password must not be empty"})
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(f.Password), 12)
	if err != nil {
		return c.JSON(500, JsonResponse{Error: err})
	}

	u := models.User{
		Username: f.Username,
		Nickname: f.Nick,
		Password: string(hashedPw),
	}

	var count int64
	c.DB.Find(&models.User{}, "username = ?", f.Username).Count(&count)
	if count > 0 {
		return c.JSON(409, JsonResponse{Message: "User with that username already exists"})
	}

	if tx := c.DB.Create(&u); tx.Error != nil {
		return c.JSON(500, JsonResponse{Error: err})
	}

	return c.JSON(200, JsonResponse{Success: true, Message: "Registration was successful"})
}
