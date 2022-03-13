package models

type User struct {
	Common
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"-"`
}
