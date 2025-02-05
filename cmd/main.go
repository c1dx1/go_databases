package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/config"
	"main/internal/databases"
	"main/internal/repositories"
	"main/internal/tasks"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %v", err))
	}

	db, err := databases.ConnectDB(cfg.PostgresURL())
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %v", err))
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("error closing database connection: %v", err))
		}
	}(db)

	gorm, err := databases.ConnectGORM(cfg.GORMURL())
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database gorm: %v", err))
	}

	pgxpool, err := databases.ConnectDBPool(cfg.PostgresURL())
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database pgxpool: %v", err))
	}
	defer pgxpool.Close()

	defaultRepo := repositories.NewDefaultDBRepository(db)
	gormRepo := repositories.NewGormRepository(gorm)
	pgxpoolRepo := repositories.NewPgxPoolRepository(pgxpool)

	task := tasks.NewTasks(defaultRepo, gormRepo, pgxpoolRepo)

	err = task.SelectTask()
	if err != nil {
		log.Fatal(fmt.Errorf("error in tasks: %v", err))
	}
}
