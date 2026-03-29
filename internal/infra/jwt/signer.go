package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
}

func NewSigner(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, keyID string) *Signer {
	if keyID == "" {
		keyID = "gofiber-hax-key-1"
	}
	return &Signer{privateKey: privateKey, publicKey: publicKey, keyID: keyID}
}

// ----------------------------------------------------
// 1. ฟังก์ชันสร้าง/เซ็นต์ JWT
// ----------------------------------------------------
func (s *Signer) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.keyID
	return token.SignedString(s.privateKey)
}

// ----------------------------------------------------
// 2. ฟังก์ชันแปลงร่างกุญแจเป็น JWKS Public Key
// ----------------------------------------------------
func (s *Signer) GetJWKSet() map[string]any {
	n := base64.RawURLEncoding.EncodeToString(s.publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.publicKey.E)).Bytes())
	return map[string]any{
		"keys": []map[string]any{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": s.keyID,
				"n":   n,
				"e":   e,
			},
		},
	}
}

func (s *Signer) PublicKey() crypto.PublicKey {
	return s.publicKey
}

// ----------------------------------------------------
// 3. ฟังก์ชันสร้าง/อ่านไฟล์กุญแจ (ตัวช่วยของแพกเกจ jwt)
// ----------------------------------------------------
func LoadKeys(privateKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, nil, errors.New("invalid PEM private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

func LoadOrGenerateKeys(privateKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, pub, err := LoadKeys(privateKeyPath)
	if err == nil {
		return priv, pub, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, nil, fmt.Errorf("load private key: %w", err)
	}

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	_ = os.MkdirAll(filepath.Dir(privateKeyPath), 0700)
	_ = os.WriteFile(privateKeyPath, privPEM, 0600)
	return priv, &priv.PublicKey, nil
}
