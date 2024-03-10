package swagger

import (
	"time"
)

type GetHistoryMatrix struct {
	TimeStart  time.Time `json:"time_start"`
	TimeEnd    time.Time `json:"time_end"`
	IsBaseline bool      `json:"is_baseline,omitempty"`
}
