package mongo

import (
	"testing"

	d "gofiber-hax/internal/core/domain"
)

func TestMongoUserMappingRoundTrip(t *testing.T) {
	input := d.Users{
		AccountID: "0001",
		Fname:     "John",
		Lname:     "Doe",
		FullName:  "John Doe",
		Username:  "jdoe",
		Password:  "hashed",
		Email:     "john@example.com",
		Phone:     "0812345678",
	}

	model := ToMongoUser(input)
	output := ToDomainUser(&model)

	if output.AccountID != input.AccountID || output.Email != input.Email || output.Password != input.Password {
		t.Fatalf("round trip mismatch: got %+v want %+v", output, input)
	}
}
