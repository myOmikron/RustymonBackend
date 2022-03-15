package handler

import (
	"RustymonBackend/models"
	"RustymonBackend/utils"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"time"
)

type LoginForm struct {
	Username string
	Password string
}

func Login(c utils.Context) error {
	var l LoginForm
	b, _ := ioutil.ReadAll(c.Request().Body)

	if err := json.Unmarshal(b, &l); err != nil {
		return c.JSON(400, JsonResponse{Error: err})
	}

	if l.Username == "" || l.Password == "" {
		return c.JSON(400, JsonResponse{Message: "Parameter username and password must not be empty"})
	}

	var u models.User
	var count int64

	c.DB.Find(&u, "username = ?", l.Username).Count(&count)
	if count == 0 {
		bcrypt.CompareHashAndPassword([]byte("hash ™"), []byte("Password to deny time based enumeration of users"))
		return c.JSON(401, JsonResponse{Message: "Login failed"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password)); err != nil {
		return c.JSON(401, JsonResponse{Message: "Login failed"})
	}

	var sessionKey string
	for {
		secureBytes := make([]byte, 64)
		if _, err := rand.Read(secureBytes); err != nil {
			continue
		}
		sessionKey = fmt.Sprintf("%x", secureBytes)
		var count int64
		c.DB.Find(&models.Session{}, "session_key = ?", sessionKey).Count(&count)
		if count == 0 {
			break
		}
	}

	s := models.Session{
		UserID:     u.ID,
		SessionKey: sessionKey,
		ValidUntil: time.Now().UTC().Add(24 * time.Hour), // 1 day valid
	}
	c.DB.Create(&s)

	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   sessionKey,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
		MaxAge:  int((time.Hour * 24).Seconds()),
	})
	return c.JSON(200, JsonResponse{Success: true})
}
