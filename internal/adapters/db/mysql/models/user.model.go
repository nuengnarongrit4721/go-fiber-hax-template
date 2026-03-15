package models

type Users struct {
	BaseModel
	AccountID string `gorm:"column:account_id;type:varchar(64);index"`
	Fname     string `gorm:"column:fname;type:varchar(100)"`
	Lname     string `gorm:"column:lname;type:varchar(100)"`
	FullName  string `gorm:"column:full_name;type:varchar(200)"`
	Username  string `gorm:"column:username;type:varchar(100)"`
	Password  string `gorm:"column:password;type:varchar(255)"`
	Email     string `gorm:"column:email;type:varchar(191);uniqueIndex"`
	Phone     string `gorm:"column:phone;type:varchar(30)"`
}

func (Users) TableName() string {
	return "users"
}
