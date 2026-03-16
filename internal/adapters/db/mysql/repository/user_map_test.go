package repository

import (
	"testing"
	"time"

	d "gofiber-hax/internal/core/domain"

	"github.com/google/uuid"
)

func TestUserMappingRoundTrip(t *testing.T) {
	now := time.Now().UTC()
	id := uuid.New()

	input := d.Users{
		BaseDomain: d.BaseDomain{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		AccountID: "0001",
		Fname:     "John",
		Lname:     "Doe",
		FullName:  "John Doe",
		Username:  "jdoe",
		Password:  "hashed",
		Email:     "john@example.com",
		Phone:     "0812345678",
	}

	model := ToModelUser(&input)
	output := ToDomainUser(&model)

	if output.ID != input.ID || output.Email != input.Email || output.Password != input.Password {
		t.Fatalf("round trip mismatch: got %+v want %+v", output, input)
	}
}
