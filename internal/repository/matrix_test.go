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

func TestMatrixRepoGetHistory(t *testing.T) {
	db := initDB()

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := time.Now().Add(time.Hour * 24)

	data := models.GetHistoryMatrix{
		TimeStart:  yesterday,
		TimeEnd:    tomorrow,
		IsBaseline: null.NewBool(false, false),
	}

	repo := InitMatrixRepo(db, 100)
	res, err := repo.GetHistory(context.Background(), data)
	fmt.Println(res)
	fmt.Println(err)
}

func TestMatrixRepoCreateGetDifference(t *testing.T) {
	db := initDB()

	repo := InitMatrixRepo(db, 100)

	data1 := models.MatrixBase{
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
	data2 := models.MatrixBase{
		Name:       "File_2",
		IsBaseLine: false,
		ParentName: null.NewString("File_1", true),
		Data: []models.MatrixNode{
			{
				MicroCategoryID: 1,
				RegionID:        1,
				Price:           100,
			},
			{
				MicroCategoryID: 2,
				RegionID:        2,
				Price:           1500,
			},
			{
				MicroCategoryID: 3,
				RegionID:        4,
				Price:           230,
			},
		},
	}

	name, err := repo.CreateMatrix(context.Background(), data1)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)

	name, err = repo.CreateMatrix(context.Background(), data2)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)

	difference, err := repo.GetDifference(context.Background(), "File_1", "File_2")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", difference)
}

func TestMatrixRepo_GetPriceTendency(t *testing.T) {
	db := initDB()

	repo := InitMatrixRepo(db, 100)

	// from postgres time to go time.Time
	timestampStr := "2024-03-10 19:52:37.053174"
	layout := "2006-01-02 15:04:05.000000"
	timestamp, err := time.Parse(layout, timestampStr)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return
	}

	tomorrow := time.Now().Add(time.Hour * 24)

	data := models.GetTendencyNode{
		MicrocategoryID: 2,
		RegionID:        2,
		TimeStart:       timestamp,
		TimeEnd:         tomorrow,
	}

	fmt.Println(timestamp)
	fmt.Println(tomorrow)

	tendency, err := repo.GetPriceTendency(context.Background(), data)

	fmt.Println(tendency)
	fmt.Println(err)
}
