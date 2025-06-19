package pkg

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

// ParseRSAPrivateKey принимает срез байтов с данными PEM и возвращает RSA приватный ключ.
func ParseRSAPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	pemData = []byte(strings.TrimSpace(string(pemData)))

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid private key PEM: decoding failed")
	}

	if block.Type == "PRIVATE KEY" {
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		if rsaPriv, ok := priv.(*rsa.PrivateKey); ok {
			return rsaPriv, nil
		}

		return nil, errors.New("not an RSA private key")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM type")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ParseRSAPublicKey принимает срез байтов с данными PEM и возвращает RSA публичный ключ.
func ParseRSAPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	pemData = []byte(strings.TrimSpace(string(pemData)))

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid public key PEM: decoding failed")
	}

	if block.Type == "PUBLIC KEY" {
		pubIfc, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		pub, ok := pubIfc.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}

		return pub, nil
	}

	return nil, errors.New("invalid public key PEM type")
}
