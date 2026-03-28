package dto

type RegisterRequest struct {
	Fname           string `json:"fname" validate:"required,min=2"`
	Lname           string `json:"lname" validate:"required,min=2"`
	Username        string `json:"username" validate:"required,min=4"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required,min=10"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}
