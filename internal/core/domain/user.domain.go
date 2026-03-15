package domain

type Users struct {
	BaseDomain
	AccountID string `json:"account_id"`
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
