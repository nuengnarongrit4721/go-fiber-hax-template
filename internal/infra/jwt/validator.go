package jwt

import (
	"errors"
	"fmt"
	"strings"

	"gofiber-hax/internal/infra/config"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	cfg  config.JWTConfig
	jwks *JwksCache
	algs []string
}

func NewValidator(cfg config.JWTConfig) *Validator {
	algs := cfg.AllowedAlgs
	if len(algs) == 0 {
		algs = []string{"RS256"}
	}
	return &Validator{
		cfg:  cfg,
		jwks: NewJwksCache(cfg.JWKSURL, cfg.CacheTTL),
		algs: algs,
	}
}

func (v *Validator) Validate(token string) (gojwt.MapClaims, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	claims := gojwt.MapClaims{}
	parser := gojwt.NewParser(
		gojwt.WithValidMethods(v.algs),
		gojwt.WithLeeway(v.cfg.ClockSkew),
	)

	parsed, err := parser.ParseWithClaims(token, claims, func(t *gojwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		return v.jwks.GetKey(kid)
	})
	if err != nil {
		return nil, err
	}
	if parsed == nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	issuer, _ := claims.GetIssuer()
	if !issuerAllowed(issuer, v.cfg.Issuer) {
		return nil, fmt.Errorf("invalid issuer: %s", issuer)
	}

	audience, _ := claims.GetAudience()
	if v.cfg.Audience != "" && !audienceAllowed(audience, v.cfg.Audience) {
		return nil, fmt.Errorf("invalid audience")
	}

	return claims, nil
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
		if (a == "https://accounts.google.com" && issuer == "accounts.google.com") ||
			(a == "accounts.google.com" && issuer == "https://accounts.google.com") {
			return true
		}
	}

	return false
}

func audienceAllowed(aud gojwt.ClaimStrings, required string) bool {
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
