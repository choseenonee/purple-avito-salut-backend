package models

import (
	"github.com/guregu/null"
	"time"
)

type MatrixNode struct {
	MicroCategoryID int `json:"microcategory_id"`
	RegionID        int `json:"region_id"`
	Price           int `json:"price"`
}

// с фронта без UNIX, надо будет приклеить
type MatrixBase struct {
	Name       string       `json:"name"`
	IsBaseLine bool         `json:"is_baseline"`
	ParentName null.String  `json:"parent_name"`
	Data       []MatrixNode `json:"data"`
}

// на бэке вставляем таймстамп
type Matrix struct {
	MatrixBase
	TimeStamp time.Time `json:"timestamp"`
}
