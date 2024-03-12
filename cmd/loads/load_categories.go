package main

import (
	"context"
	"fmt"
	"template/pkg/config"
	"template/pkg/database"
)

func main() {
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

	config.InitConfig()
	db := database.GetDB()
	for _, i := range data {
		_, err := db.ExecContext(context.Background(), `INSERT INTO microcategories (name) VALUES ($1);`, i)
		if err != nil {
			panic(err.Error())
		}

	}
	fmt.Println("okay!")

	newData := map[int][]int{
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

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}

	for i := 1; i < 25; i++ {
		for _, val := range newData[i] {
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
