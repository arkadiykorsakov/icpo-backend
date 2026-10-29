package city

import (
	"time"

	countryDomain "icpo-backend/internal/domain/country"
)

type City struct {
	ID        int64                 `db:"id"            json:"id"`
	Name      string                `db:"name"          json:"name"`
	Country   countryDomain.Country `db:"country"       json:"country"`
	CreatedAt time.Time             `db:"created_at"    json:"created_at"`
}
