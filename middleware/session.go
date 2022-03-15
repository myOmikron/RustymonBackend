package middleware

import (
	"RustymonBackend/middleware/session"
	"RustymonBackend/models"
	"RustymonBackend/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Session(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var auth bool
			var userID uint
			var s models.Session
			if cookie, err := c.Request().Cookie("session_id"); errors.Is(err, http.ErrNoCookie) {
				auth = false
				userID = 0
			} else {
				var count int64
				db.Find(&s, "session_key = ?", cookie.Value).Count(&count)
				if count == 0 {
					auth = false
					userID = 0
				} else {
					if !s.ValidUntil.UTC().After(time.Now().UTC()) {
						auth = false
						userID = 0
					} else {
						auth = true
						userID = s.UserID
					}
				}
			}

			cc := &utils.Context{
				Context: c,
				DB:      db,
				Session: session.Create(auth, userID),
			}
			return next(cc)
		}
	}
}
