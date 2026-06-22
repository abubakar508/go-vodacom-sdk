package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// EncryptPKCS1v15ToBase64 encrypts value with an RSA public key and returns a
// base64 string suitable for the M-Pesa Authorization Bearer value.
//
// The official PHP implementation uses openssl_public_encrypt, whose default
// padding is RSA PKCS#1 v1.5. M-Pesa OpenAPI expects that same encryption mode.
func EncryptPKCS1v15ToBase64(value, publicKey string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("value cannot be empty")
	}

	key, err := ParseRSAPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(value))
	if err != nil {
		return "", fmt.Errorf("encrypt rsa pkcs1v15: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// ParseRSAPublicKey accepts either a raw base64 DER public key, a PEM public
// key, a PKIX SubjectPublicKeyInfo key, or a PKCS#1 RSA public key.
func ParseRSAPublicKey(publicKey string) (*rsa.PublicKey, error) {
	clean := strings.TrimSpace(publicKey)
	if clean == "" {
		return nil, errors.New("public key cannot be empty")
	}

	var der []byte
	if block, _ := pem.Decode([]byte(clean)); block != nil {
		der = block.Bytes
	} else {
		clean = strings.ReplaceAll(clean, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		clean = strings.ReplaceAll(clean, " ", "")
		decoded, err := base64.StdEncoding.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf("decode public key: %w", err)
		}
		der = decoded
	}

	parsed, err := x509.ParsePKIXPublicKey(der)
	if err == nil {
		key, ok := parsed.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is %T, not *rsa.PublicKey", parsed)
		}
		return key, nil
	}

	key, pkcs1Err := x509.ParsePKCS1PublicKey(der)
	if pkcs1Err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("parse rsa public key: pkix=%v pkcs1=%v", err, pkcs1Err)
}
