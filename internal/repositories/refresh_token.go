package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at,
			replaced_by_token_id
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.ReplacedByTokenID,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at,
			replaced_by_token_id,
			created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	token := &models.RefreshToken{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.ReplacedByTokenID,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int64) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE id = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RefreshTokenRepository) SetReplacement(ctx context.Context, id int64, replacedByTokenID int64) error {
	query := `
		UPDATE refresh_tokens
		SET replaced_by_token_id = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		replacedByTokenID,
		id,
	)

	return err
}
