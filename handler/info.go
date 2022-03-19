package handler

import (
	u "github.com/myOmikron/echotools/utility"
)

type ServerInfo struct {
	Version uint16 `json:"version"`
}

func Info(c *Context) error {
	serverInfo := ServerInfo{
		Version: 1,
	}

	return c.JSON(200, u.JsonResponse{Success: true, Data: &serverInfo})
}
