package models

type Users struct {
	BaseModel
	AccountID string `json:"account_id" bson:"account_id"`
	Fname     string `json:"fname" bson:"fname"`
	Lname     string `json:"lname" bson:"lname"`
	FullName  string `json:"full_name" bson:"full_name"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
	Phone     string `json:"phone" bson:"phone"`
}
