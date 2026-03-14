package models

type Users struct {
	ID    string `gorm:"column:id;primaryKey"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func (Users) TableName() string {
	return "users"
}
