package models

type UserBase struct {
	Name     string `json:"name"`
	RegionID int    `json:"region_id"`
}
