package service

import (
	"context"
	"errors"
	"testing"

	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/out"
	errs "gofiber-hax/internal/shared/errors"
)

type stubUserRepo struct {
	createFn func(ctx context.Context, req *d.Users) error
}

func (s stubUserRepo) CreateUser(ctx context.Context, req *d.Users) error {
	if s.createFn != nil {
		return s.createFn(ctx, req)
	}
	return nil
}

func (s stubUserRepo) GetByAccountID(ctx context.Context, accountID string) (d.Users, error) {
	return d.Users{}, nil
}

func (s stubUserRepo) GetByUsername(ctx context.Context, username string) (d.Users, error) {
	return d.Users{}, nil
}

var _ out.UserRepository = stubUserRepo{}

func TestCreateUserServiceGeneratesAccountIDAndNormalizesEmail(t *testing.T) {
	var captured d.Users
	svc := NewUserService(stubUserRepo{
		createFn: func(ctx context.Context, req *d.Users) error {
			captured = *req
			return nil
		},
	})

	req := &d.Users{
		Email:    "  USER@Example.COM  ",
		Username: " demo ",
	}
	if err := svc.CreateUserService(context.Background(), req); err != nil {
		t.Fatalf("CreateUserService() error = %v", err)
	}

	if captured.AccountID == "" {
		t.Fatal("expected generated account id")
	}
	if captured.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", captured.Email)
	}
	if captured.Username != "demo" {
		t.Fatalf("expected trimmed username, got %q", captured.Username)
	}
}

func TestCreateUserServiceRejectsNilRequest(t *testing.T) {
	svc := NewUserService(stubUserRepo{})

	err := svc.CreateUserService(context.Background(), nil)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
