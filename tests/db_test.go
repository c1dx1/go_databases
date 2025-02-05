package tests

import (
	"main/internal/models"
	"main/internal/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSelectSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewDefaultDBRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "city"}).
		AddRow(1, "Alice", "New York").
		AddRow(2, "Bob", "Los Angeles")

	mock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)

	users, err := repo.SimpleSelectSQL()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectSQLWithParam(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewDefaultDBRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "city"}).
		AddRow(1, "Alice", "New York")

	mock.ExpectQuery("SELECT \\* FROM users WHERE 1=1 AND name ILIKE \\$1").
		WithArgs("%Alice%").
		WillReturnRows(rows)

	filters := map[string]interface{}{"name": "%Alice%"}
	users, err := repo.SelectSQLWithParam(filters)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Alice", users[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertUserSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewDefaultDBRepository(db)
	user := &models.User{Name: "Charlie", City: "Chicago"}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Name, user.City).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	err = repo.InsertUserSQL(user)
	assert.NoError(t, err)
	assert.Equal(t, 3, user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewDefaultDBRepository(db)
	user := models.User{ID: 1, Name: "UpdatedName", City: "UpdatedCity"}

	mock.ExpectExec("UPDATE users SET name = \\$1, city = \\$2 WHERE id = \\$3").
		WithArgs(user.Name, user.City, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUserSQL(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionUserSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewDefaultDBRepository(db)
	user := &models.User{Name: "New User", City: "New City"}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users WHERE id=\\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Name, user.City).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	err = repo.TransactionUserSQL(1, user)
	assert.NoError(t, err)
	assert.Equal(t, 2, user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
