package dto

type UpdateProfileRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,notblank,max=100"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,notblank"`
	NewPassword string `json:"new_password" validate:"required,min=8,maxbytes=72,passwordchars"`
}
