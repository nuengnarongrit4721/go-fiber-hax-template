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

type LoginUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
