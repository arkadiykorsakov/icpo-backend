package dto

type CreateRequest struct {
	Name      string `json:"name" validate:"required,notblank,max=255"`
	CountryID int64  `json:"country_id" validate:"required,gt=0"`
}
