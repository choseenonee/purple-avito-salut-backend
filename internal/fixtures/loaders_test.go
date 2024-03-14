package fixtures

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

func initDB() *sqlx.DB {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		"purple", "purple123", "localhost", "5432", "purple_hack",
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	}

	return db
}

func TestLoadMicroCategories(t *testing.T) {
	data := []string{
		"ROOT",
		"Техника",
		"Бытовая",
		"Компьютеры",
		"Ноутбуки",
		"Комплектующие",
		"Игровые",
		"Рабочие",
		"Процессоры",
		"Видеокарты",
		"Одежда",
		"Зима",
		"Лето",
		"Демисезонная",
		"Куртка",
		"Пальто",
		"Термобелье",
		"Мужская",
		"Женская",
		"Шорты",
		"Футболка",
		"Обувь",
		"Аксессуары",
		"Игрушки",
		"Пистолет",
		"Пони",
	}

	db := initDB()

	for _, i := range data {
		_, err := db.ExecContext(context.Background(), `INSERT INTO microcategories (name) VALUES ($1);`, i)
		if err != nil {
			panic(err.Error())
		}

		//idRaw, err := res.LastInsertId()
		//if err != nil {
		//	panic(err.Error())
		//}

		//fmt.Println(idRaw, i)
	}
	fmt.Println("okay!")
}

func TestMicroCategoriesRelationsLoad(t *testing.T) {
	data := map[int][]int{
		1:  {2, 11, 24},
		2:  {3, 4},
		4:  {5, 6},
		5:  {7, 8},
		6:  {9, 10},
		11: {12, 13, 14},
		12: {15, 16, 17},
		13: {18, 19},
		14: {22, 23},
		18: {20, 21},
		24: {25, 26},
	}

	db := initDB()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}

	for i := 1; i < 25; i++ {
		for _, val := range data[i] {
			_, err = tx.ExecContext(context.Background(), `INSERT INTO relationships_microcategories (parent_id, child_id) VALUES ($1, $2);`,
				i, val)
			if err != nil {
				_ = tx.Rollback()
				panic(err.Error())
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("okay!")
}

func TestLoadLocations(t *testing.T) {
	data := []string{
		"ROOT",
		"Якутия",
		"Саха",
		"Якутск",
		"Крым",
		"Севастополь",
		"Симферополь",
		"село Рузаевка",
		"село Авангард",
		"Ярославль",
		"Ярославский край",
		"Казань",
		"Ленинский",
		"Северный",
		"село Рыбинск",
	}

	db := initDB()

	for _, i := range data {
		g, err := db.ExecContext(context.Background(), `INSERT INTO regions (name) VALUES ($1);`, i)
		if err != nil {
			panic(err.Error())
		}

		_ = g

		//idRaw, err := res.LastInsertId()
		//if err != nil {
		//	panic(err.Error())
		//}

		//fmt.Println(idRaw, i)
	}
	fmt.Println("okay!")
}

func TestMicroLocationsRelationsLoad(t *testing.T) {
	data := map[int][]int{
		1:  {2, 5, 11},
		2:  {3, 4},
		5:  {6, 7},
		7:  {8, 9},
		11: {10, 12},
		12: {13, 14},
		13: {15},
	}

	db := initDB()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}

	for i := 1; i < 25; i++ {
		for _, val := range data[i] {
			_, err = tx.ExecContext(context.Background(), `INSERT INTO relationships_regions (parent_id, child_id) VALUES ($1, $2);`,
				i, val)
			if err != nil {
				_ = tx.Rollback()
				panic(err.Error())
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("okay!")
}

func TestBaseLineMatrixLoad(t *testing.T) {
	// cat loc price
	data := [][3]int{
		{1, 1, 1},
		{2, 2, 1},
		{2, 4, 1},
		{2, 5, 1},
		{11, 15, 1},
		{24, 4, 1},
		{13, 11, 1},
		{13, 2, 1},
		{6, 4, 1},
		{9, 2, 1},
		{21, 11, 1},
		{21, 9, 1},
	}

	db := initDB()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < len(data); i++ {
		_, err = tx.ExecContext(context.Background(), `INSERT INTO matrix (name, microcategory_id, region_id, price) VALUES ('baseline_test', $1, $2, $3);`,
			data[i][0], data[i][1], data[i][2])
		if err != nil {
			_ = tx.Rollback()
			panic(err.Error())
		}
	}

	_, err = tx.ExecContext(context.Background(), `INSERT INTO matrix_metadata 
    (matrix_name, timestamp, is_baseline, parent_matrix_name) VALUES ('baseline_test', $1, 'true', $2);`, time.Now(), nil)

	err = tx.Commit()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("okay!")
}

func TestDiscountMatrixLoad(t *testing.T) {
	// cat loc price
	datas := [][][3]int{
		{
			{5, 12, 10000},
			{5, 7, 10001},
			{14, 14, 10002},
		},
		{
			{4, 2, 200},
			{4, 13, 200},
			{11, 13, 200},
			{11, 2, 200},
		},
	}

	db := initDB()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}
	for d := 0; d < len(datas); d++ {
		for i := 0; i < len(datas[d]); i++ {
			_, err = tx.ExecContext(context.Background(), `INSERT INTO matrix (name, microcategory_id, region_id, price) VALUES ($4, $1, $2, $3);`,
				datas[d][i][0], datas[d][i][1], datas[d][i][2], fmt.Sprintf("discount_%v", d))
			if err != nil {
				_ = tx.Rollback()
				panic(err.Error())
			}
		}

		_, err = tx.ExecContext(context.Background(), `INSERT INTO matrix_metadata 
    	(matrix_name, timestamp, is_baseline, parent_matrix_name) VALUES ($3, $1, 'false', $2);`, time.Now(), nil, fmt.Sprintf("discount_%v", d))
	}

	err = tx.Commit()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("okay!")
}
