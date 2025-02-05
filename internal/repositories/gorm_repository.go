package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"main/internal/models"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

//crud - операции чисто для примера

// create
func (r *GormRepository) UserCreate() {
	usersCreate := []models.User{
		{Name: "Максим", City: "Moscow"},
		{Name: "Екатерина", City: "Novgorod"},
	}
	r.db.Create(&usersCreate)
}

// read
func (r *GormRepository) UserRead() {
	var usersRead []models.User
	r.db.Find(&usersRead)

	for _, user := range usersRead {
		fmt.Printf("ID: %d, Name: %s, Age: %s\n", user.ID, user.Name, user.City)
	}
}

// read with filters
func (r *GormRepository) UserReadWithFilters() {
	var usersRead []models.User
	r.db.Find(&usersRead, "name = ?", "Максим")

	for _, user := range usersRead {
		fmt.Printf("ID: %d, Name: %s, Age: %s\n", user.ID, user.Name, user.City)
	}
}

// update
func (r *GormRepository) UserUpdate() {
	user := models.User{
		ID: 7,
	}

	r.db.Model(&user).Update("name", "Misha")
}

// delete
func (r *GormRepository) UserDelete() {
	user := models.User{
		ID: 7,
	}

	r.db.Delete(&user)
}
