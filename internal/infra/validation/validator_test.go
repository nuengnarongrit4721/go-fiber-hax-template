package validation

import (
	"errors"
	"testing"

	errs "gofiber-hax/internal/shared/errors"
)

type sampleRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := ValidateStruct(&sampleRequest{
			Email:           "user@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := ValidateStruct(&sampleRequest{
			Email:           "bad-email",
			Password:        "short",
			ConfirmPassword: "mismatch",
		})
		if !errors.Is(err, errs.ErrInvalidInput) {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}

		var validationErr Error
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected validation error, got %T", err)
		}
		if len(validationErr.Fields) < 3 {
			t.Fatalf("expected field errors, got %+v", validationErr.Fields)
		}
	})
}
