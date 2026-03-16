package service

import (
	"context"
	"fmt"

	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/in"
	errs "gofiber-hax/internal/shared/errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService in.UserService
}

func NewAuthService(userService in.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) RegisterService(ctx context.Context, req *d.RegisterUserInput) error {
	if req == nil {
		return fmt.Errorf("authservice.register error: %w", errs.ErrInvalidInput)
	}
	if err := validateRegisterInput(req); err != nil {
		return err
	}
	if req.Password != req.ConfirmPassword {
		return fmt.Errorf("authservice.register error: %w", errs.ErrInvalidInput)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("authservice.register error: failed to hash password")
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

func validateRegisterInput(req *d.RegisterUserInput) error {
	required := map[string]string{
		"fname":            req.Fname,
		"lname":            req.Lname,
		"username":         req.Username,
		"email":            req.Email,
		"phone":            req.Phone,
		"password":         req.Password,
		"confirm_password": req.ConfirmPassword,
	}
	for _, value := range required {
		if value == "" {
			return fmt.Errorf("authservice.register error: %w", errs.ErrInvalidInput)
		}
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
