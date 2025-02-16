package models

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
	Coins    int64  `json:"coins" db:"coins"`
}
