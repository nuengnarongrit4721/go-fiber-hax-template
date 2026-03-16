package domain

type RegisterUserInput struct {
	Fname           string
	Lname           string
	Username        string
	Email           string
	Phone           string
	Password        string
	ConfirmPassword string
}
