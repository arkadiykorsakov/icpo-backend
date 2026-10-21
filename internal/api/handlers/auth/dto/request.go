package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,notblank,max=100"`
	Password string `json:"password" validate:"required,min=8,maxbytes=72,passwordchars"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
