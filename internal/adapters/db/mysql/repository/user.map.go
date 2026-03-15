package repository

import (
	m "gofiber-hax/internal/adapters/db/mysql/models"
	d "gofiber-hax/internal/core/domain"
)

func ToDomainUser(doc *m.Users) d.Users {
	if doc == nil {
		return d.Users{}
	}
	return d.Users{
		BaseDomain: d.BaseDomain{
			ID:        doc.ID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			DeletedAt: doc.DeletedAt,
		},
		AccountID: doc.AccountID,
		Fname:     doc.Fname,
		Lname:     doc.Lname,
		FullName:  doc.FullName,
		Username:  doc.Username,
		Password:  doc.Password,
		Email:     doc.Email,
		Phone:     doc.Phone,
	}
}

func ToModelUser(doc *d.Users) m.Users {
	if doc == nil {
		return m.Users{}
	}
	return m.Users{
		BaseModel: m.BaseModel{
			ID:        doc.ID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			DeletedAt: doc.DeletedAt,
		},
		AccountID: doc.AccountID,
		Fname:     doc.Fname,
		Lname:     doc.Lname,
		FullName:  doc.FullName,
		Username:  doc.Username,
		Password:  doc.Password,
		Email:     doc.Email,
		Phone:     doc.Phone,
	}
}
