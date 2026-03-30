package app

import (
	"context"
	"path/filepath"
	"testing"

	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/out"
	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/jwt"
)

type stubUserRepo struct{}

func (stubUserRepo) CreateUser(ctx context.Context, req *d.Users) error {
	return nil
}

func (stubUserRepo) GetByAccountID(ctx context.Context, accountID string) (d.Users, error) {
	return d.Users{}, nil
}

func (stubUserRepo) GetByUsername(ctx context.Context, username string) (d.Users, error) {
	return d.Users{}, nil
}

var _ out.UserRepository = stubUserRepo{}

func TestBuildServicesSkipsAuthWhenSignerIsNil(t *testing.T) {
	services := buildServices(config.Config{}, Repos{User: stubUserRepo{}}, nil)

	if services.User == nil {
		t.Fatal("expected user service")
	}
	if services.Auth != nil {
		t.Fatal("expected auth service to be nil when signer is nil")
	}
}

func TestBuildHandlersSkipsAuthRoutesWhenSignerIsNil(t *testing.T) {
	services := buildServices(config.Config{}, Repos{User: stubUserRepo{}}, nil)
	set := buildHandlers(services, nil, nil)

	if set.HTTP.V1.Auth != nil {
		t.Fatal("expected auth handler to be nil when signer is nil")
	}
	if set.HTTP.V1.JWKS != nil {
		t.Fatal("expected jwks handler to be nil when signer is nil")
	}
}

func TestBuildHandlersWiresAuthWhenSignerExists(t *testing.T) {
	priv, pub, err := jwt.LoadOrGenerateKeys(filepath.Join(t.TempDir(), "jwt_private.pem"))
	if err != nil {
		t.Fatalf("LoadOrGenerateKeys() error = %v", err)
	}
	signer := jwt.NewSigner(priv, pub, "test-key")

	services := buildServices(config.Config{}, Repos{User: stubUserRepo{}}, signer)
	set := buildHandlers(services, nil, signer)

	if set.HTTP.V1.Auth == nil {
		t.Fatal("expected auth handler when signer exists")
	}
	if set.HTTP.V1.JWKS == nil {
		t.Fatal("expected jwks handler when signer exists")
	}
}
