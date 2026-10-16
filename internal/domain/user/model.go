package user

import "time"

type User struct {
	ID           int64     `db:"id"           json:"id"`
	Email        string    `db:"email"        json:"email"`
	Name         string    `db:"name"         json:"name"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"   json:"updated_at"`
}
