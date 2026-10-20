package user

import (
	"strings"
	"time"
)

type User struct {
	ID           int64     `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	Name         string    `db:"name"          json:"name"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

func NormalizeEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}

func NewUser(email, name string) User {
	email = NormalizeEmail(email)
	name = strings.TrimSpace(name)

	return User{
		Name:  name,
		Email: email,
	}
}
