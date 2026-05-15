package storage

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
)

// signJWT creates a signed JWT using the service account's private key
func (d *DriveStorage) signJWT(claims map[string]interface{}) (string, error) {
	// Parse private key
	block, _ := pem.Decode([]byte(d.credentials.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block from private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not RSA")
	}

	// Create JWT header
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerB64 := base64URLEncode(headerJSON)
	claimsB64 := base64URLEncode(claimsJSON)

	signingInput := headerB64 + "." + claimsB64

	// Sign
	hashed := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	signatureB64 := base64URLEncode(signature)

	return signingInput + "." + signatureB64, nil
}

func base64URLEncode(data []byte) string {
	s := base64.StdEncoding.EncodeToString(data)
	s = strings.TrimRight(s, "=")
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	return s
}
