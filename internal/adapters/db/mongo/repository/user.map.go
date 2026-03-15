package mongo

import (
	m "gofiber-hax/internal/adapters/db/mongo/models"
	d "gofiber-hax/internal/core/domain"
)

func ToDomainUser(doc *m.Users) d.Users {
	if doc == nil {
		return d.Users{}
	}
	return d.Users{
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

func ToMongoUser(domain d.Users) m.Users {
	return m.Users{
		AccountID: domain.AccountID,
		Fname:     domain.Fname,
		Lname:     domain.Lname,
		FullName:  domain.FullName,
		Username:  domain.Username,
		Password:  domain.Password,
		Email:     domain.Email,
		Phone:     domain.Phone,
	}
}
