package service

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"template/internal/repository"
	"testing"
)

func initDB() *sqlx.DB {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		"postgres", "postgres", "localhost", "5432", "postgres",
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	}

	return db
}

func TestUpdate(t *testing.T) {
	db := initDB()

	repo := repository.InitMatrixRepo(db, 100)

	service := InitUpdateService(repo)

	preparedStorage, err := service.PrepareStorage(context.Background(), "baseline_test", []string{"discount_0", "discount_1"})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(preparedStorage.MicroCategoryHops)
	fmt.Println(preparedStorage.RegionHops)
	fmt.Println(preparedStorage.DiscountHops)
}
