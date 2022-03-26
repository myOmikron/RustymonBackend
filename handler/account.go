package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/RustymonBackend/tasks"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/middleware"
	u "github.com/myOmikron/echotools/utility"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
	"net/mail"
)

var (
	ErrLoginFailed   = errors.New("login failed")
	ErrUsernameTaken = errors.New("username is taken")
	ErrEmailTaken    = errors.New("email is taken")
)

type AccountHandler struct {
	DB         *gorm.DB
	Config     *configs.RustymonConfig
	WorkerPool worker.Pool
}

type RegisterForm struct {
	Username    *string `json:"username" echotools:"required;not empty"`
	Password    *string `json:"password" echotools:"required;not empty"`
	Email       *string `json:"email" echotools:"required;not empty"`
	TrainerName *string `json:"trainer_name" echotools:"required;not empty"`
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

		if address, err := mail.ParseAddress(*f.Email); err != nil {
			return c.JSON(400, u.JsonResponse{Error: fmt.Sprintf("No valid mail provided: %s", err.Error())})
		} else {
			*f.Email = address.Address
		}

		var userCount, emailCount int64
		a.DB.Find(&utilitymodels.User{}, "username = ?", *f.Username).Count(&userCount)
		a.DB.Find(&utilitymodels.User{}, "email = ?", *f.Email).Count(&emailCount)
		if emailCount > 0 {
			return c.JSON(409, u.JsonResponse{Error: ErrEmailTaken.Error()})
		}
		if userCount > 0 {
			return c.JSON(409, u.JsonResponse{Error: ErrUsernameTaken.Error()})
		}

		user, err := database.CreateUser(a.DB, *f.Username, *f.Password, *f.Email, true)
		if err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}
		player := models.Player{
			User:        *user,
			TrainerName: *f.TrainerName,
		}
		if err = a.DB.Create(&player).Error; err != nil {
			return c.JSON(500, u.JsonResponse{Error: err.Error()})
		}

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}

type LoginForm struct {
	Username *string `json:"username" echotools:"required;not empty"`
	Password *string `json:"password" echotools:"required;not empty"`
}

func (a *AccountHandler) Login() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		var l LoginForm
		if err := u.ValidateJsonForm(c, &l); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		user, err := auth.Authenticate(a.DB, *l.Username, *l.Password)
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

type ResetPasswordUsernameForm struct {
	Username *string `json:"username" echotools:"required;not empty"`
}

func (a *AccountHandler) ResetPasswordUsername() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		var f ResetPasswordUsernameForm

		if err := u.ValidateJsonForm(c, &f); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		var user utilitymodels.User
		var userCount int64
		a.DB.Where("username = ?", *f.Username).Find(&user).Count(&userCount)
		if userCount == 1 {
			mailTask := tasks.NewMailTask(&a.Config.Mail, []string{user.Email}, "Password reset", "Test mail")
			a.WorkerPool.AddTask(mailTask)
		} else {
			c.Logger().Info("Found multiple matching users")
		}

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}

type ResetPasswordEmailForm struct {
	Email *string `json:"email" echotools:"required;not empty"`
}

func (a *AccountHandler) ResetPasswordEmail() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		var f ResetPasswordEmailForm

		if err := u.ValidateJsonForm(c, &f); err != nil {
			return c.JSON(400, u.JsonResponse{Error: err.Error()})
		}

		var userCount int64
		var user utilitymodels.User
		a.DB.Where("email = ?", *f.Email).Find(&user).Count(&userCount)
		if userCount == 1 {
			mailTask := tasks.NewMailTask(&a.Config.Mail, []string{user.Email}, "Password reset", "Test mail")
			a.WorkerPool.AddTask(mailTask)
		} else {
			c.Logger().Infof("Found %d matching users", userCount)
		}

		return c.JSON(200, u.JsonResponse{Success: true})
	})
}
