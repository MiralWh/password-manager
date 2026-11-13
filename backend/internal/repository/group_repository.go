package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/models"
)

type GroupRepository struct {
	pool *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{pool: pool}
}

func (r *GroupRepository) List(ctx context.Context, userID int64, search string) ([]models.Group, error) {
	query := `SELECT g.id, g.user_id, g.name, g.created_at,
	                 COUNT(p.id) FILTER (WHERE p.id IS NOT NULL) AS cnt
	          FROM groups g
	          LEFT JOIN passwords p ON p.group_id = g.id
	          WHERE g.user_id=$1`
	args := []interface{}{userID}
	if search != "" {
		query += ` AND g.name ILIKE '%'||$2||'%'`
		args = append(args, search)
	}
	query += ` GROUP BY g.id ORDER BY g.name`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt, &g.PasswordCount); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *GroupRepository) Create(ctx context.Context, userID int64, name string) (*models.Group, error) {
	var g models.Group
	err := r.pool.QueryRow(ctx,
		`INSERT INTO groups (user_id, name) VALUES ($1,$2) RETURNING id, user_id, name, created_at`,
		userID, name,
	).Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt)
	return &g, err
}

func (r *GroupRepository) Rename(ctx context.Context, userID, id int64, name string) (*models.Group, error) {
	var g models.Group
	err := r.pool.QueryRow(ctx,
		`UPDATE groups SET name=$1 WHERE id=$2 AND user_id=$3 RETURNING id, user_id, name, created_at`,
		name, id, userID,
	).Scan(&g.ID, &g.UserID, &g.Name, &g.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &g, err
}

func (r *GroupRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM groups WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
