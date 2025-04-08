package pgrepo

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgOauth struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgOauth(db *postgres.Storage, log *zap.Logger) *pgOauth {
	return &pgOauth{
		log: log,
		db:  db,
	}
}

func (p *pgOauth) LinkAccount(
	ctx context.Context,
	clientID string,
	provider *models.OAuthProvider,
) error {
	const op = "repository.pg.OAuth.LinkAccount"

	query, args, err := p.db.Builder.
		Insert("oauth_accounts").
		Columns("client_id", "provider", "provider_id", "access_token", "expires_at").
		Values(clientID, provider.Provider, provider.ProviderID, provider.AccessToken, provider.ExpiresAt).
		ToSql()
	if err != nil {
		p.log.Error(op+" failed to build SQL query", zap.Error(err))
		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+" failed to execute SQL query", zap.Error(err))
		return err
	}

	return nil
}

func (p *pgOauth) GetByProvider(
	ctx context.Context,
	provider string,
	providerID string,
) (*models.OAuthProvider, error) {
	const op = "repository.pg.OAuth.GetByProvider"

	query, args, err := p.db.Builder.
		Select("client_id", "provider", "provider_id", "access_token", "expires_at").
		From("oauth_accounts").
		Where(squirrel.Eq{"provider": provider, "provider_id": providerID}).
		ToSql()
	if err != nil {
		p.log.Error(op+" failed to build SQL query", zap.Error(err))
		return nil, err
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	var oauthProvider models.OAuthProvider
	err = row.Scan(
		&oauthProvider.ClientID,
		&oauthProvider.Provider,
		&oauthProvider.ProviderID,
		&oauthProvider.AccessToken,
		&oauthProvider.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Info(op+" OAuth provider not found", zap.String("provider", provider), zap.String("providerID", providerID))
			return nil, domain.ErrNotFound
		}
		p.log.Error(op+" failed to scan row", zap.Error(err))
		return nil, err
	}

	return &oauthProvider, nil
}

func (p *pgOauth) UnlinkAccount(ctx context.Context, clientID string, provider string) error {
	const op = "repository.pg.OAuth.UnlinkAccount"

	query, args, err := p.db.Builder.
		Delete("oauth_accounts").
		Where(squirrel.Eq{"client_id": clientID, "provider": provider}).
		ToSql()
	if err != nil {
		p.log.Error(op+" failed to build SQL query", zap.Error(err))
		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+" failed to execute SQL query", zap.Error(err))
		return err
	}

	return nil
}
