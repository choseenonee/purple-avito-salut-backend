package models

import "github.com/guregu/null"

type RegionBase struct {
	ParentID null.Int
	Name     string
}

type UserSegmentBase struct {
	UserID    int `json:"user_id"`
	SegmentID int `json:"segment_id"`
}
