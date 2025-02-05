package databases

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
