package models

type Storage struct {
	BaseLineMatrix   Matrix   `json:"baseline"`
	DiscountMatrices []Matrix `json:"discount"`
}

type PreparedStorage struct {
	Storage
	MicroCategoryHops []int
	RegionHops        []int
	DiscountHops      map[string]map[int]map[int]int
}
