package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/middleware"
	u "github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/gorm"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

var (
	ErrLoginFailed   = errors.New("login failed")
	ErrUsernameTaken = errors.New("username is taken")
	ErrEmailTaken    = errors.New("email is taken")
)

type AccountHandler struct {
	DB     *gorm.DB
	Config *configs.RustymonConfig
}

type RegisterForm struct {
	Username    string `json:"username" echotools:"required;not empty"`
	Password    string `json:"password" echotools:"required;not empty"`
	Email       string `json:"email" echotools:"required;not empty"`
	TrainerName string `json:"trainer_name" echotools:"required;not empty"`
}

func (a *AccountHandler) Register() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		if a.Config.Rustymon.RegistrationDisabled {
			return c.JSON(503, u.JsonResponse{Error: "registration is disabled on this server"})
		}
		var f RegisterForm

		if err := u.ValidateJsonForm(c, &f); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		if address, err := mail.ParseAddress(f.Email); err != nil {
			return c.JSON(400, u.JsonResponse{Error: fmt.Sprintf("No valid mail provided: %s", err.Error())})
		} else {
			f.Email = address.Address
		}

		var userCount, emailCount int64
		a.DB.Find(&utilitymodels.User{}, "username = ?", f.Username).Count(&userCount)
		a.DB.Find(&utilitymodels.User{}, "email = ?", f.Email).Count(&emailCount)
		if emailCount > 0 {
			return c.JSON(409, u.JsonResponse{Error: ErrEmailTaken.Error()})
		}
		if userCount > 0 {
			return c.JSON(409, u.JsonResponse{Error: ErrUsernameTaken.Error()})
		}

		user, err := database.CreateUser(a.DB, f.Username, f.Password, f.Email, true)
		if err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}
		player := models.Player{
			User:        *user,
			TrainerName: f.TrainerName,
		}
		if err = a.DB.Create(&player).Error; err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}

type LoginForm struct {
	Username string `json:"username" echotools:"required;not empty"`
	Password string `json:"password" echotools:"required;not empty"`
}

func (a *AccountHandler) Login() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		var l LoginForm
		if err := u.ValidateJsonForm(c, &l); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		user, err := auth.Authenticate(a.DB, l.Username, l.Password)
		if err != nil || user == nil {
			c.Logger().Info(err)
			return c.JSON(403, u.JsonResponse{Error: ErrLoginFailed.Error()})
		} else {
			if err := middleware.Login(a.DB, user, c); err != nil {
				return c.JSON(500, u.JsonResponse{Error: err.Error()})
			}
		}

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}

func (a *AccountHandler) Logout() echo.HandlerFunc {
	return middleware.LoginRequired(func(c *Context) error {
		if err := middleware.Logout(a.DB, c); err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}
		return c.JSON(200, u.JsonResponse{Success: true})
	})
}

type ResetPasswordForm struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

func (a *AccountHandler) ResetPassword() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		var f ResetPasswordForm

		if err := u.ValidateJsonForm(c, &f); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		if f.Email == nil && f.Username == nil {
			return c.JSON(400, u.JsonResponse{Error: "username or email is required"})
		}

		if *f.Email == "" && *f.Username == "" {
			return c.JSON(400, u.JsonResponse{Error: "username or email must not be empty"})
		}

		// Receiver email address.
		to := []string{
			"crsi@hopfen.space",
		}
		from := a.Config.Mail.User

		// Message.
		message := []byte(fmt.Sprintf(`To: %s
From: %s
Subject: Very important information regarding your Rustymon subscription™

The Rustymon backend is now able to utilize the power of MAIL.

Enjoy!
`, strings.Join(to, ", "), from))

		// Authentication.
		au := smtp.PlainAuth("", from, a.Config.Mail.Password, a.Config.Mail.Host)

		// Sending email.
		err := smtp.SendMail(net.JoinHostPort(a.Config.Mail.Host, strconv.Itoa(int(a.Config.Mail.Port))), au, from, to, message)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, u.JsonResponse{Error: "error while sending mail"})
		}
		fmt.Println("Email Sent Successfully!")

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}
