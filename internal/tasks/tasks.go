package tasks

import (
	"fmt"
	"main/internal/models"
	"main/internal/repositories"
	"main/migrator"
)

type Tasks struct {
	dbRepository      *repositories.DefaultDBRepository
	gormRepository    *repositories.GormRepository
	pgxpoolRepository *repositories.PgxPoolRepository
}

func NewTasks(dbRepo *repositories.DefaultDBRepository, gormRepo *repositories.GormRepository, pgxpoolRepo *repositories.PgxPoolRepository) *Tasks {
	return &Tasks{dbRepository: dbRepo, gormRepository: gormRepo, pgxpoolRepository: pgxpoolRepo}
}

func (t *Tasks) SelectTask() error {
	for {
		taskNumber := 0
		fmt.Println("\nВыберите номер задания:\n" +
			"2 - простой SELECT-запрос для выборки данных из таблицы\n" +
			"3 - SELECT-запрос с параметром\n4 - Вставка новой записи в таблицу\n5 - Обновление строки в таблице\n" +
			"6 - Транзакция для добавления новой записи и удалении старой\n7 - Тест GORM\n8 - Запуск миграций\n" +
			"9 - Создать пул соединений и выполнить тестовые запросы\n0 - Выйти")
		_, err := fmt.Scan(&taskNumber)
		if err != nil {
			return fmt.Errorf("tasks.go: SelectTask: fmt.Scan: %w", err)
		}
		if taskNumber == 0 {
			return nil
		}

		switch taskNumber {
		case 2:
			err = t.SimpleSelectSQL()
			if err != nil {
				return err
			}
		case 3:
			err = t.SelectSQLWithParam()
			if err != nil {
				return err
			}
		case 4:
			err = t.InsertUserSQL()
			if err != nil {
				return err
			}
		case 5:
			err = t.UpdateUserSQL()
			if err != nil {
				return err
			}
		case 6:
			err = t.TransactionSQL()
			if err != nil {
				return err
			}
		case 7:
			t.GORMTest()
		case 8:
			err = t.RunMigrations()
			if err != nil {
				return err
			}
		case 9:
			err = t.DBPools()
			if err != nil {
				return err
			}
		default:
			fmt.Println("Incorrect task number")
		}
	}
}

// zadanie 2
func (t *Tasks) SimpleSelectSQL() error {
	users, err := t.dbRepository.SimpleSelectSQL()
	if err != nil {
		return fmt.Errorf("tasks.go: SimpleSelectSQL: queries.SimpleSQLQuery: %s", err)
	}
	fmt.Println(users)
	return nil
}

// zadanie 3
func (t *Tasks) SelectSQLWithParam() error {
	filters := make(map[string]interface{})

	var user models.User
	fmt.Println("Введите ID для фильтрации (если не нужно фильтровать по ID - введите 0):")
	_, err := fmt.Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("tasks.go: SelectSQLWithParam: fmt.Scan ID: %s", err)
	}
	//здесь это корректно, так как первичный ключ начинает отсчет от 1
	if user.ID != 0 {
		filters["id"] = user.ID
	}

	fmt.Println("Введите имя для фильтрации (если не нужно фильтровать по имени - просто нажмите Enter):")
	fmt.Scanln(&user.Name)
	if user.Name != "" {
		filters["name"] = user.Name
	}

	fmt.Println("Введите город для фильтрации (если не нужно фильтровать по городу - просто нажмите Enter):")
	_, err = fmt.Scanln(&user.City)
	if user.City != "" {
		filters["city"] = user.City
	}

	users, err := t.dbRepository.SelectSQLWithParam(filters)
	if err != nil {
		return fmt.Errorf("tasks.go: SelectSQLWithParam: queries.SelectSQLWithParam: %s", err)
	}
	fmt.Println(users)
	return nil
}

// zadanie 4
func (t *Tasks) InsertUserSQL() error {
	var user models.User

	fmt.Println("Введите имя нового user (оно не должно быть пустым):")
	_, err := fmt.Scan(&user.Name)
	if err != nil {
		return fmt.Errorf("tasks.go: InsertUserSQL: fmt.Scan Name: %s", err)
	}
	if user.Name == "" {
		return fmt.Errorf("tasks.go: InsertUserSQL: fmt.Scan Name: %s", "user name is empty")
	}

	fmt.Println("Введите город нового user:")
	fmt.Scanln(&user.City)

	err = t.dbRepository.InsertUserSQL(&user)
	if err != nil {
		return fmt.Errorf("tasks.go: InsertUserSQL: queries.InsertUserSQL: %s", err)
	}

	fmt.Printf("New user ID: %d", user.ID)

	return nil
}

// zadanie 5
// здесь нет возможности изменить на пустое значение, так как fmt.Scan не позволяет сканить значения в модель с указателями
func (t *Tasks) UpdateUserSQL() error {
	var user models.User
	nulls := 0

	fmt.Println("Введите ID user'а, данные которого хотите изменить: ")
	_, err := fmt.Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("tasks.go: UpdateUserSQL: fmt.Scan ID: %s", err)
	}

	if user.ID <= 0 {
		return fmt.Errorf("tasks.go: UpdateUserSQL: fmt.Scan ID: %s", "user id is empty")
	}

	fmt.Println("Введите новое имя user'а (если имя менять не нужно, то просто нажмите Enter): ")
	fmt.Scanln(&user.Name)

	if user.Name == "" {
		nulls++
	}

	fmt.Println("Введите новый город user'а (если город менять не нужно, то просто нажмите Enter): ")
	fmt.Scanln(&user.City)

	if user.City == "" {
		nulls++
	}

	if nulls == 2 {
		return fmt.Errorf("tasks.go: UpdateUserSQL: empty updates")
	}

	err = t.dbRepository.UpdateUserSQL(user)
	if err != nil {
		return fmt.Errorf("tasks.go: UpdateUserSQL: queries.UpdateUserSQL: %s", err)
	}
	fmt.Println("Success")
	return nil
}

// zadanie 6
func (t *Tasks) TransactionSQL() error {
	var delID int

	fmt.Println("Введите ID сотрудника, которого заменяют")
	_, err := fmt.Scan(&delID)
	if err != nil {
		return fmt.Errorf("tasks.go: TransactionSQL: fmt.Scan ID: %s", err)
	}

	var user models.User

	fmt.Println("Введите имя сотрудника, которым заменяют предыдущего (оно не должно быть пустым):")
	_, err = fmt.Scan(&user.Name)
	if err != nil {
		return fmt.Errorf("tasks.go: TransactionSQL: fmt.Scan Name: %s", err)
	}
	if user.Name == "" {
		return fmt.Errorf("tasks.go: TransactionSQL: fmt.Scan Name: %s", "user name is empty")
	}

	fmt.Println("Введите город сотрудника, которым заменяют предыдущего:")
	fmt.Scanln(&user.City)

	err = t.dbRepository.TransactionUserSQL(delID, &user)
	if err != nil {
		return fmt.Errorf("tasks.go: TransactionSQL: queries.TransactionUserSQL: %s", err)
	}
	fmt.Printf("New user ID: %d", user.ID)

	return nil
}

// zadanie 7
func (t *Tasks) GORMTest() {
	t.gormRepository.UserCreate()

	t.gormRepository.UserRead()

	t.gormRepository.UserReadWithFilters()

	t.gormRepository.UserUpdate()

	t.gormRepository.UserDelete()
}

// zadanie 8
func (t *Tasks) RunMigrations() error {
	err := migrator.RunMigrations()
	if err != nil {
		return fmt.Errorf("tasks.go: RunMigrations: %s", err)
	}

	return nil
}

// zadanie 9
func (t *Tasks) DBPools() error {
	users, err := t.pgxpoolRepository.SelectSQL()
	if err != nil {
		return fmt.Errorf("tasks.go: DBPools: SelectSQL: %s", err)
	}

	fmt.Println(users)
	return nil
}
