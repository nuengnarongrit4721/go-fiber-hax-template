package models

type User struct {
	BaseModel
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func (User) TableName() string {
	return "users"
}
