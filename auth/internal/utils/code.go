package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
)

const (
	minCode  = 9000
	aligment = 1000
)

func GenerateSecureCode() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(minCode))
	if err != nil {
		return 0, fmt.Errorf("%w %w", domain.ErrCodeConfirmGen, err)
	}

	return int(n.Int64()) + aligment, nil
}
