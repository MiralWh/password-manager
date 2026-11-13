package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/models"
)

type PasswordRepository struct {
	pool *pgxpool.Pool
}

func NewPasswordRepository(pool *pgxpool.Pool) *PasswordRepository {
	return &PasswordRepository{pool: pool}
}

type PasswordRow struct {
	models.Password
	EncPassword []byte
}

const cols = `id, user_id, group_id, title, username, encrypted_password, url, notes, created_at, updated_at`

func scanRow(row pgx.Row) (*PasswordRow, error) {
	var p PasswordRow
	err := row.Scan(&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Username, &p.EncPassword,
		&p.URL, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PasswordRepository) ListByUser(ctx context.Context, userID int64, search string) ([]PasswordRow, error) {
	query := `SELECT p.id, p.user_id, p.group_id, p.title, p.username, p.encrypted_password,
	                 p.url, p.notes, p.created_at, p.updated_at
	          FROM passwords p
	          LEFT JOIN groups g ON g.id = p.group_id
	          WHERE p.user_id=$1`
	args := []interface{}{userID}
	if search != "" {
		query += ` AND (
			p.title    ILIKE '%'||$2||'%' OR
			p.username ILIKE '%'||$2||'%' OR
			p.url      ILIKE '%'||$2||'%' OR
			p.notes    ILIKE '%'||$2||'%' OR
			COALESCE(g.name, '') ILIKE '%'||$2||'%'
		)`
		args = append(args, search)
	}
	query += ` ORDER BY p.updated_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []PasswordRow
	for rows.Next() {
		var p PasswordRow
		if err := rows.Scan(&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Username, &p.EncPassword,
			&p.URL, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *PasswordRepository) Get(ctx context.Context, userID, id int64) (*PasswordRow, error) {
	p, err := scanRow(r.pool.QueryRow(ctx,
		`SELECT `+cols+` FROM passwords WHERE id=$1 AND user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (r *PasswordRepository) Create(ctx context.Context, userID int64, in models.PasswordInput, enc []byte) (*PasswordRow, error) {
	return scanRow(r.pool.QueryRow(ctx,
		`INSERT INTO passwords (user_id, group_id, title, username, encrypted_password, url, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING `+cols,
		userID, in.GroupID, in.Title, in.Username, enc, in.URL, in.Notes))
}

func (r *PasswordRepository) Update(ctx context.Context, userID, id int64, in models.PasswordInput, enc []byte) (*PasswordRow, error) {
	p, err := scanRow(r.pool.QueryRow(ctx,
		`UPDATE passwords SET group_id=$1, title=$2, username=$3, encrypted_password=$4,
		        url=$5, notes=$6, updated_at=NOW()
		 WHERE id=$7 AND user_id=$8 RETURNING `+cols,
		in.GroupID, in.Title, in.Username, enc, in.URL, in.Notes, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (r *PasswordRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM passwords WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
