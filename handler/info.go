package handler

import "RustymonBackend/utils"

type ServerInfo struct {
	Version uint16 `json:"version"`
}

func Info(c utils.Context) error {
	serverInfo := ServerInfo{
		Version: 1,
	}
	return c.JSON(200, JsonResponse{Success: true, Data: &serverInfo})
}
