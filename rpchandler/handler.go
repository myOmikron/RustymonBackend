package rpchandler

import (
	"errors"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/logging"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/gorm"
	"net/mail"
)

var (
	ErrInvalidEmail   = errors.New("invalid email")
	ErrUsernameTaken  = errors.New("username taken")
	ErrEmailTaken     = errors.New("email taken")
	ErrDatabaseError  = errors.New("database error")
	ErrParameterEmpty = errors.New("parameter is empty")
)

var log = logging.GetLogger("rpc-handler")

type RPC struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Username    string
	Email       string
	Password    string
	TrainerName string
}

type RegisterResult struct {
	ErrorMessage string
}

func (r *RPC) RegisterUser(req RegisterRequest, res *RegisterResult) error {
	if req.Username == "" || req.Email == "" || req.Password == "" || req.TrainerName == "" {
		return ErrParameterEmpty
	}

	if address, err := mail.ParseAddress(req.Email); err != nil {
		res.ErrorMessage = err.Error()
		return ErrInvalidEmail
	} else {
		req.Email = address.Address
	}

	var userCount, emailCount, confirmEmailCount int64
	r.DB.Find(&utilitymodels.User{}, "username = ?", req.Username).Count(&userCount)
	r.DB.Find(&utilitymodels.User{}, "email = ?", req.Email).Count(&emailCount)
	r.DB.Find(&models.PlayerConfirmEmail{}, "email = ?", req.Email).Count(&confirmEmailCount)
	if emailCount > 0 || confirmEmailCount > 0 {
		return ErrEmailTaken
	}
	if userCount > 0 {
		return ErrUsernameTaken
	}

	user, err := database.CreateUser(r.DB, req.Username, req.Password, &req.Email, true)
	if err != nil {
		res.ErrorMessage = err.Error()
		return ErrDatabaseError
	}
	player := models.Player{
		User:        *user,
		TrainerName: req.TrainerName,
	}
	if err = r.DB.Create(&player).Error; err != nil {
		res.ErrorMessage = err.Error()
		return ErrDatabaseError
	}

	log.Infof("Registered %s per CLI tool", req.Username)

	return nil
}
