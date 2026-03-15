package service

import (
	"context"
	"errors"
	"fmt"
	"gofiber-hax/internal/adapters/http/dto"
	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/in"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService in.UserService
}

func NewAuthService(userService in.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) RegisterService(ctx context.Context, req *dto.RegisterRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("password and confirm password do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	req.Password = string(hashedPassword)

	newUser := d.Users{
		Fname:    req.Fname,
		Lname:    req.Lname,
		FullName: req.Fname + " " + req.Lname,
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		Phone:    req.Phone,
	}

	if err := s.userService.CreateUserService(ctx, &newUser); err != nil {
		return fmt.Errorf("authservice.register error: %w", err)
	}

	return nil
}

func (s *AuthService) LoginService() {

}

func (s *AuthService) LogoutService() {

}

func (s *AuthService) RefreshTokenService() {

}

func (s *AuthService) ForgotPasswordService() {

}

func (s *AuthService) ResetPasswordService() {

}

func (s *AuthService) ChangePasswordService() {

}

func (s *AuthService) VerifyEmailService() {

}

func (s *AuthService) ResendVerificationEmailService() {

}
