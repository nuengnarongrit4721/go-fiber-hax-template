package models

type Users struct {
	BaseModel
	AccountID string `gorm:"column:account_id:index"`
	Fname     string `gorm:"column:fname"`
	Lname     string `gorm:"column:lname"`
	FullName  string `gorm:"column:full_name"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	Email     string `gorm:"column:email:index;unique"`
	Phone     string `gorm:"column:phone"`
}

func (Users) TableName() string {
	return "users"
}
