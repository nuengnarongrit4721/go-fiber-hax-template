package handlers

import (
	"gofiber-hax/internal/adapters/http/dto"
	d "gofiber-hax/internal/core/domain"
)

func ToUserResponse(domainUser d.Users) dto.UserResponse {
	return dto.UserResponse{
		ID:        domainUser.ID.String(),
		AccountID: domainUser.AccountID,
		FirstName: domainUser.Fname,
		LastName:  domainUser.Lname,
		Username:  domainUser.Username,
		Email:     domainUser.Email,
		Phone:     domainUser.Phone,
	}
}
