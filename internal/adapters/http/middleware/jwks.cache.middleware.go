package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type jwksCache struct {
	url       string
	ttl       time.Duration
	client    *http.Client
	mu        sync.Mutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
}

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func newJWKSCache(url string, ttl time.Duration) *jwksCache {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &jwksCache{
		url: url,
		ttl: ttl,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		keys: make(map[string]*rsa.PublicKey),
	}
}

func (c *jwksCache) GetKey(kid string) (*rsa.PublicKey, error) {
	if kid == "" {
		return nil, errors.New("missing kid")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if key, ok := c.keys[kid]; ok && time.Since(c.fetchedAt) < c.ttl {
		return key, nil
	}

	if err := c.refreshLocked(); err != nil {
		return nil, err
	}

	if key, ok := c.keys[kid]; ok {
		return key, nil
	}

	return nil, fmt.Errorf("kid not found: %s", kid)
}

func (c *jwksCache) refreshLocked() error {
	if c.url == "" {
		return errors.New("jwks url is empty")
	}

	resp, err := c.client.Get(c.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks fetch failed: %s", resp.Status)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, k := range jwks.Keys {
		if k.Kty != "RSA" || k.Kid == "" {
			continue
		}
		pub, err := jwkToPublicKey(k)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}

	if len(keys) == 0 {
		return errors.New("jwks contains no usable keys")
	}

	c.keys = keys
	c.fetchedAt = time.Now()
	return nil
}

func jwkToPublicKey(k jwkKey) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	if e.Sign() <= 0 || e.BitLen() > 31 {
		return nil, errors.New("invalid exponent")
	}
	eInt := int(e.Int64())
	if eInt <= 0 || eInt > math.MaxInt32 {
		return nil, errors.New("invalid exponent range")
	}

	return &rsa.PublicKey{N: n, E: eInt}, nil
}
