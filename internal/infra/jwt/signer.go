package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewSigner(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *Signer {
	return &Signer{privateKey: privateKey, publicKey: publicKey}
}

// ----------------------------------------------------
// 1. ฟังก์ชันสร้าง/เซ็นต์ JWT
// ----------------------------------------------------
func (s *Signer) Sign(claims jwt.MapClaims) (string, error) {
	// สร้างตัว Token และระบุว่าใช้ Algorithm เป็น RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// สำคัญมาก! ต้องระบุรหัสกุญแจ (Key ID) เพื่อให้ฝั่ง Validator รู้ว่าต้องใช้อะไรไข
	token.Header["kid"] = "gofiber-hax-key-1"
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
				"kid": "gofiber-hax-key-1",
				"n":   n,
				"e":   e,
			},
		},
	}
}

// ----------------------------------------------------
// 3. ฟังก์ชันสร้าง/อ่านไฟล์กุญแจ (ตัวช่วยของแพกเกจ jwt)
// ----------------------------------------------------
func LoadOrGenerateKeys(privateKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err == nil {
		block, _ := pem.Decode(keyBytes)
		if block != nil {
			priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err == nil {
				return priv, &priv.PublicKey, nil
			}
		}
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
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
