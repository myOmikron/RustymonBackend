package handler

import (
	"github.com/myOmikron/RustymonBackend/configs"
	u "github.com/myOmikron/echotools/utility"
)

type ServerInfo struct {
	Version              uint16 `json:"version"`
	RegistrationDisabled bool   `json:"registration_disabled"`
}

func Info(c *Context) error {
	serverInfo := ServerInfo{
		Version:              1,
		RegistrationDisabled: configs.Config.Rustymon.RegistrationDisabled,
	}

	return c.JSON(200, u.JsonResponse{Success: true, Data: &serverInfo})
}
