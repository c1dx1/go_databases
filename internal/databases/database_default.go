package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// zadanie 1
func ConnectDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("postgres.go: ConnectDB: sql.open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres.go: ConnectDB: db.ping: %v", err)
	}

	return db, nil
}
