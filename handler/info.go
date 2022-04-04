package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/echotools/middleware"
	u "github.com/myOmikron/echotools/utility"
	"gorm.io/gorm"
)

type ServerInfoHandler struct {
	Config *configs.RustymonConfig
	DB     *gorm.DB
}

type serverInfoData struct {
	Version               uint16 `json:"version"`
	RegistrationDisabled  bool   `json:"registration_disabled"`
	PasswordResetDisabled bool   `json:"password_reset_disabled"`
}

func (s *ServerInfoHandler) Serverinfo() echo.HandlerFunc {
	return middleware.Wrap(func(c *Context) error {
		serverInfo := serverInfoData{
			Version:              1,
			RegistrationDisabled: s.Config.Rustymon.RegistrationDisabled,
		}

		return c.JSON(200, u.JsonResponse{Success: true, Data: &serverInfo})
	})
}
