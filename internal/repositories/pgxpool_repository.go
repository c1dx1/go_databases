package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"main/internal/models"
)

type PgxPoolRepository struct {
	pool *pgxpool.Pool
}

func NewPgxPoolRepository(pool *pgxpool.Pool) *PgxPoolRepository {
	return &PgxPoolRepository{pool: pool}
}

func (r *PgxPoolRepository) SelectSQL() ([]models.User, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pgxpool_repository.go: SelectSQL: r.pool.Acquire(): %w", err)
	}
	defer conn.Release()

	rows, err := r.pool.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("pgxpool_repository.go: SelectSQL: r.pool.Query(): %w", err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.City)
		if err != nil {
			return nil, fmt.Errorf("pgxpool_repository.go: SelectSQL: rows.Scan(): %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
