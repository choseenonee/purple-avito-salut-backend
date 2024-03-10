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

type GetDifferenceByMatricesNames struct {
	Name1 string `json:"from_name"`
	Name2 string `json:"to_name"`
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

type ResponseHistoryMatrix struct {
	Name       string      `json:"name"`
	TimeStamp  time.Time   `json:"timestamp"`
	ParentName null.String `json:"parent_name"`
}

type GetTendencyNode struct {
	MicrocategoryID int       `json:"microcategory_id"`
	RegionID        int       `json:"region_id"`
	TimeStart       time.Time `json:"time_start"`
	TimeEnd         time.Time `json:"time_end"`
}

type ResponseTendencyNode struct {
	TimeStamp time.Time `json:"timestamp"`
	Price     int       `json:"price"`
}
