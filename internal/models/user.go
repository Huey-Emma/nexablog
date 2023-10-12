package models

import (
	"time"

	"nexablog/pkg/lib"
	"nexablog/pkg/validator"
)

type User struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymousUser() bool {
	return u == AnonymousUser
}

func (u *User) IsOwner(identifier int) bool {
	if u.IsAnonymousUser() {
		return false
	}

	if u.UserID == identifier {
		return true
	}

	return false
}

type UserIn struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUser(v *validator.Validator, u *UserIn) {
	v.Check(lib.NonWhiteSpace(u.Username), "username", "cannot be blank")
	v.Check(lib.NonWhiteSpace(u.Email), "email", "cannot be blank")
	v.Check(lib.NonWhiteSpace(u.Password), "password", "cannot be blank")
	v.Check(lib.ValidEmail(u.Email), "email", "provide a valid email")
	v.Check(len(u.Password) >= 8, "password", "must be at least 8 characters")
}
