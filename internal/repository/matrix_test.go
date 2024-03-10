package repository

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"testing"
	"time"
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

func TestMatrixRepo_GetHistory(t *testing.T) {
	db := initDB()

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := time.Now().Add(time.Hour * 24)

	data := models.GetHistoryMatrix{
		TimeStart:  yesterday,
		TimeEnd:    tomorrow,
		IsBaseline: null.NewBool(false, false),
	}

	repo := InitMatrixRepo(db)
	res, err := repo.GetHistory(context.Background(), data)
	fmt.Println(res)
	fmt.Println(err)
}

func TestMatrixRepo_Create(t *testing.T) {
	db := initDB()

	repo := InitMatrixRepo(db)
	data := models.MatrixBase{
		Name:       "File_1",
		IsBaseLine: false,
		ParentName: null.NewString("", false),
		Data: []models.MatrixNode{
			{
				MicroCategoryID: 1,
				RegionID:        1,
				Price:           100,
			},
			{
				MicroCategoryID: 2,
				RegionID:        2,
				Price:           200,
			},
			{
				MicroCategoryID: 3,
				RegionID:        3,
				Price:           300,
			},
		},
	}

	name, err := repo.CreateMatrix(context.Background(), data)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func TestMatrixRepo_GetPriceTendency(t *testing.T) {
	db := initDB()

	repo := InitMatrixRepo(db)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := time.Now().Add(time.Hour * 24)

	data := models.GetTendencyNode{
		MicrocategoryID: 1,
		RegionID:        1,
		TimeStart:       yesterday,
		TimeEnd:         tomorrow,
	}

	tendency, err := repo.GetPriceTendency(context.Background(), data)

	fmt.Println(tendency)
	fmt.Println(err)
}
