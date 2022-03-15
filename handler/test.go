package handler

import (
	"RustymonBackend/utils"
)

func Test(c utils.Context) error {
	return c.JSON(200, JsonResponse{Success: true, Data: c.Session.IsAuthenticated()})
}
