package pgrepo

import (
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgToken struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgToken(db *postgres.Storage, log *zap.Logger) *pgToken {
	return &pgToken{
		log: log,
		db:  db,
	}
}
