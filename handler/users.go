package handler

import (
	"RustymonBackend/models"
	"RustymonBackend/utils"
	"net/http"
)

func GetUsers(c utils.Context) error {
	var u []models.User
	c.DB.Find(&u)
	return c.JSON(http.StatusOK, &JsonResponse{
		Success: true,
		Data:    u,
	})
}
