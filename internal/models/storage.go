package models

type StorageBase struct {
	BaseLineMatrixName  string   `json:"baseline"`
	DiscountMatrixNames []string `json:"discount"`
}

type Storage struct {
	BaseLineMatrix   Matrix   `json:"baseline"`
	DiscountMatrices []Matrix `json:"discount"`
}

type PreparedStorage struct {
	Storage
	MicroCategoryHops []int                          `json:"micro_category_hops"`
	RegionHops        []int                          `json:"region_hops"`
	DiscountHops      map[string]map[int]map[int]int `json:"discount_hops"`
	SegmentDiscount   map[int]string                 `json:"segment_discount"`
}

type PreparedStorageSend struct {
	StorageBase
	MicroCategoryHops []int                          `json:"micro_category_hops"`
	RegionHops        []int                          `json:"region_hops"`
	DiscountHops      map[string]map[int]map[int]int `json:"discount_hops"`
	SegmentDiscount   map[int]string                 `json:"segment_discount"`
}
