package middleware

import (
	"errors"
	"fmt"
	"strings"

	"gofiber-hax/internal/infra/config"

	"github.com/golang-jwt/jwt/v5"
)

type jwtValidator struct {
	cfg  config.JWTConfig
	jwks *jwksCache
	algs []string
}

func newJWTValidator(cfg config.JWTConfig) *jwtValidator {
	algs := cfg.AllowedAlgs
	if len(algs) == 0 {
		algs = []string{"RS256"}
	}
	return &jwtValidator{
		cfg:  cfg,
		jwks: newJWKSCache(cfg.JWKSURL, cfg.CacheTTL),
		algs: algs,
	}
}

func (v *jwtValidator) Validate(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}

	claims := &jwt.RegisteredClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods(v.algs),
		jwt.WithLeeway(v.cfg.ClockSkew),
	)

	parsed, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		return v.jwks.GetKey(kid)
	})
	if err != nil {
		return err
	}
	if parsed == nil || !parsed.Valid {
		return errors.New("invalid token")
	}

	if !issuerAllowed(claims.Issuer, v.cfg.Issuer) {
		return fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}
	if v.cfg.Audience != "" && !audienceAllowed(claims.Audience, v.cfg.Audience) {
		return fmt.Errorf("invalid audience")
	}

	return nil
}

func issuerAllowed(issuer, allowed string) bool {
	issuer = strings.TrimSpace(issuer)
	if issuer == "" {
		return false
	}
	allowed = strings.TrimSpace(allowed)
	if allowed == "" {
		return true
	}

	for _, a := range strings.Split(allowed, ",") {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if issuer == a {
			return true
		}
		// Google accepts both forms
		if (a == "https://accounts.google.com" && issuer == "accounts.google.com") ||
			(a == "accounts.google.com" && issuer == "https://accounts.google.com") {
			return true
		}
	}

	return false
}

func audienceAllowed(aud jwt.ClaimStrings, required string) bool {
	if required == "" {
		return true
	}

	requiredList := strings.Split(required, ",")
	for _, req := range requiredList {
		req = strings.TrimSpace(req)
		if req == "" {
			continue
		}
		for _, got := range aud {
			if got == req {
				return true
			}
		}
	}

	return false
}
