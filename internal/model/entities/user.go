package entities

import "github.com/guregu/null"

type UserBase struct {
	// эта библиотека спасает, когда происходит попытка Scan'ировать null-value в Goшный тип данных
	Name null.String `json:"name"`
}

type UserCreate struct {
	UserBase
	Password string `json:"password"`
}

type User struct {
	UserBase
	ID int `json:"id"`
}
