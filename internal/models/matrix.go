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

type MatrixDifference struct {
	Updated [][2]MatrixNode `json:"updated"`
	Added   []MatrixNode    `json:"added"`
	Deleted []MatrixNode    `json:"deleted"`
}

// на бэке вставляем таймстамп
type Matrix struct {
	MatrixBase
	TimeStamp time.Time `json:"timestamp"`
}

type GetHistoryMatrix struct {
	TimeStart  time.Time `json:"time_start"`
	TimeEnd    time.Time `json:"time_end"`
	IsBaseline null.Bool `json:"is_baseline,omitempty"`
}
