package repository

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

func TestMatrixRepo_GetHistory(t *testing.T) {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		"postgres", "postgres", "localhost", "5432", "postgres",
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	}

	repo := InitmatrixRepo(db)
	res, err := repo.GetHistory(context.Background(), time.Now(), time.Now().Add(time.Hour*24), null.NewBool(false, false))
	fmt.Print(res)
}
