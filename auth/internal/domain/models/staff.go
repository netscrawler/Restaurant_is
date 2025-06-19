package models

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	lowercase              = "abcdefghijklmnopqrstuvwxyz"
	uppercase              = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits                 = "0123456789"
	symbols                = "!@#$%&*+-=?"
	standartPasswordLength = 8
)

type Staff struct {
	ID                 uuid.UUID // UUID
	WorkEmail          string
	PasswordHash       string
	Position           string
	IsActive           bool
	NeedChangePassword bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Staff) IsActiveString() string {
	if s.IsActive {
		return "Active"
	}

	return "Inactive"
}

func (s *Staff) UpdatePassword(newPassword string) error {
	newPswd, err := HashPassword(s.PasswordHash)
	if err != nil {
		return domain.ErrGeneratePassword
	}

	if s.NeedChangePassword {
		s.NeedChangePassword = false
	}

	s.PasswordHash = newPswd
	s.UpdatedAt = time.Now()

	return nil
}

func (s *Staff) UpdatePosition(newPosition string) {
	s.Position = newPosition
	s.UpdatedAt = time.Now()
}

func (s *Staff) UpdateEmail(newEmail string) {
	s.WorkEmail = newEmail
	s.UpdatedAt = time.Now()
}

func generatePassword(length int, useSymbols bool) (string, error) {
	charset := lowercase + uppercase + digits
	if useSymbols {
		charset += symbols
	}

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func (u *Staff) CheckPasswordHash(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func NewStaff(email, position string) (*Staff, string, error) {
	pswd, err := generatePassword(standartPasswordLength, true)
	if err != nil {
		return nil, "", fmt.Errorf("%w (%w)", domain.ErrGeneratePassword, err)
	}

	hash, err := HashPassword(pswd)
	if err != nil {
		return nil, "", fmt.Errorf("%w (%w)", domain.ErrGeneratePassword, err)
	}

	return &Staff{
		ID:                 uuid.New(),
		WorkEmail:          email,
		PasswordHash:       hash,
		Position:           position,
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		NeedChangePassword: true,
	}, pswd, nil
}
