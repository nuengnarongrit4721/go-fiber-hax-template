package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/in"
	auth "gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/jwt"
	errs "gofiber-hax/internal/shared/errors"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService in.UserService
	signer      *jwt.Signer
	cfg         auth.AuthConfig
}

func NewAuthService(
	userService in.UserService,
	signer *jwt.Signer,
	cfg auth.AuthConfig,
) *AuthService {
	return &AuthService{
		userService: userService,
		signer:      signer,
		cfg:         cfg,
	}
}

var _ in.AuthService = (*AuthService)(nil)

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
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Fname = strings.TrimSpace(req.Fname)
	req.Lname = strings.TrimSpace(req.Lname)

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

func (s *AuthService) LoginService(ctx context.Context, req *d.LoginUserInput) (string, error) {
	if req == nil {
		return "", fmt.Errorf("authservice.login error: %w", errs.ErrInvalidInput)
	}

	if err := validateLoginInput(req); err != nil {
		return "", err
	}
	req.Username = strings.TrimSpace(req.Username)

	user, err := s.userService.GetUserByUsernameService(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("authservice.login error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("authservice.login error: %w", errs.ErrUnauthorized)
	}

	now := time.Now()
	claims := gojwt.MapClaims{
		"sub":      user.AccountID,
		"jti":      uuid.NewString(),
		"user_id":  user.AccountID,
		"username": user.Username,
		"email":    user.Email,
		"iss":      s.cfg.JWT.Issuer,
		"aud":      s.cfg.JWT.Audience,
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"exp":      now.Add(s.cfg.JWT.AccessTokenTTL).Unix(),
	}

	token, err := s.signer.Sign(claims)
	if err != nil {
		return "", fmt.Errorf("authservice.login error: failed to generate token")
	}

	return token, nil
}

func validateLoginInput(req *d.LoginUserInput) error {
	required := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	for _, value := range required {
		if value == "" {
			return fmt.Errorf("authservice.login error: %w", errs.ErrInvalidInput)
		}
	}
	return nil
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
