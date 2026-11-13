package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/models"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool
}

func NewAPIKeyRepository(pool *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{pool: pool}
}

func (r *APIKeyRepository) Create(ctx context.Context, userID int64, name, keyHash, keyPrefix string) (*models.APIKey, error) {
	var k models.APIKey
	err := r.pool.QueryRow(ctx,
		`INSERT INTO api_keys (user_id, name, key_hash, key_prefix)
		 VALUES ($1,$2,$3,$4) RETURNING id, user_id, name, key_prefix, created_at, last_used_at`,
		userID, name, keyHash, keyPrefix,
	).Scan(&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.CreatedAt, &k.LastUsedAt)
	return &k, err
}

func (r *APIKeyRepository) List(ctx context.Context, userID int64) ([]models.APIKey, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, key_prefix, created_at, last_used_at
		 FROM api_keys WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.APIKey
	for rows.Next() {
		var k models.APIKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.CreatedAt, &k.LastUsedAt); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (r *APIKeyRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM api_keys WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *APIKeyRepository) FindUserByHash(ctx context.Context, keyHash string) (int64, string, int64, error) {
	var userID, keyID int64
	var username string
	err := r.pool.QueryRow(ctx,
		`SELECT ak.id, u.id, u.username
		 FROM api_keys ak JOIN users u ON u.id = ak.user_id
		 WHERE ak.key_hash=$1`, keyHash,
	).Scan(&keyID, &userID, &username)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", 0, ErrNotFound
	}
	return userID, username, keyID, err
}

func (r *APIKeyRepository) Touch(ctx context.Context, id int64) {
	_, _ = r.pool.Exec(ctx, "UPDATE api_keys SET last_used_at=NOW() WHERE id=$1", id)
}
