package repositories

import (
	"database/sql"
	"fmt"
	"main/internal/models"
)

type DefaultDBRepository struct {
	db *sql.DB
}

func NewDefaultDBRepository(db *sql.DB) *DefaultDBRepository {
	return &DefaultDBRepository{db: db}
}

func (r *DefaultDBRepository) SimpleSelectSQL() ([]models.User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("default_db_queries.go: SimpleSQLQuery: db.Query: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		user := models.User{}

		err = rows.Scan(&user.ID, &user.Name, &user.City)
		if err != nil {
			return nil, fmt.Errorf("default_db_queries.go: SimpleSQLQuery: rows.Scan: %v", err)
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *DefaultDBRepository) SelectSQLWithParam(filters map[string]interface{}) ([]models.User, error) {
	query := "SELECT * FROM users WHERE 1=1"

	args := []interface{}{}
	index := 1
	for key, value := range filters {
		if key == "id" {
			query += fmt.Sprintf(" AND %s=$%d", key, index)
			args = append(args, value)
			index++
		} else {
			query += fmt.Sprintf(" AND %s ILIKE $%d", key, index)
			args = append(args, value)
			index++
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("default_db_queries.go: SelectSQLWithParam: db.Query: %v", err)
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.City)
		if err != nil {
			return nil, fmt.Errorf("default_db_queries.go: SelectSQLWithParam: rows.Scan: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *DefaultDBRepository) InsertUserSQL(user *models.User) error {
	err := r.db.QueryRow("INSERT INTO users (name, city) VALUES ($1, $2) RETURNING id", user.Name, user.City).
		Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("default_db_queries.go: InsertUserSQL: db.Query: %v", err)
	}

	return nil
}

// обычно делаю обновление через мапу, но здесь показалось нецелесообразным
func (r *DefaultDBRepository) UpdateUserSQL(user models.User) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	index := 1
	if user.Name != "" {
		query += fmt.Sprintf("name = $%d", index)
		args = append(args, user.Name)
		index++
	}
	if user.City != "" {
		if index > 1 {
			query += ", "
		}
		query += fmt.Sprintf("city = $%d", index)
		args = append(args, user.City)
		index++
	}

	query += fmt.Sprintf(" WHERE id = $%d", index)
	args = append(args, user.ID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("default_db_queries.go: UpdateUserSQL: db.Query: %v", err)
	}

	return nil
}

func (r *DefaultDBRepository) TransactionUserSQL(delID int, user *models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("default_db_queries.go: TransactionUserSQL: db.Begin: %v", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	_, err = tx.Exec("DELETE FROM users WHERE id=$1", delID)
	if err != nil {
		return fmt.Errorf("default_db_queries.go: TransactionUserSQL: db.Exec delete: %v", err)
	}

	err = tx.QueryRow("INSERT INTO users (name, city) VALUES ($1, $2) RETURNING id", &user.Name, &user.City).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("default_db_queries.go: TransactionUserSQL: db.Exec insert: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("default_db_queries.go: TransactionUserSQL: db.Commit: %v", err)
	}
	return nil
}
