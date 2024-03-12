package service

import (
	"context"
	"template/internal/models"
	"template/internal/repository"
)

type updateService struct {
	matrixRepo repository.Matrix
}

func InitUpdateService(matrixRepo repository.Matrix) Update {
	return updateService{matrixRepo: matrixRepo}
}

// mock user segments and segment discount matrices

var userSegments = map[int]int{
	1: 100,
	2: 200,
}

var segmentMatrices = map[int]string{
	100: "discount_0",
	200: "discount_1",
}

func recursive(index int, in [][4]int, ans []int, lastWithPrice int, lastIndex int, isFirst bool) bool {
	isFound := false
	for _, i := range in[lastIndex:] {
		if isFirst && i[0] != index {
			break
		}
		if i[0] != index {
			continue
		}
		if i[2] != 0 {
			lastWithPrice = index
		}
		if !isFound {
			ans[index-1] = lastWithPrice
		}
		if !recursive(i[1], in, ans, lastWithPrice, index-1, false) {
			if i[3] != 0 {
				ans[i[1]-1] = i[1]
			} else {
				ans[i[1]-1] = lastWithPrice
			}
		}
		isFound = true
	}
	return isFound
}

func setNestedMapValue(m map[string]map[int]map[int]int, key string, k1, k2, value int) {
	// Check if the second-level map exists; if not, create it
	if m[key] == nil {
		m[key] = make(map[int]map[int]int)
	}

	// Check if the third-level map exists; if not, create it
	if m[key][k1] == nil {
		m[key][k1] = make(map[int]int)
	}

	// Now that all maps are initialized, set the value
	m[key][k1][k2] = value
}

func calculateDiscountMatrix(matrices []models.Matrix, discountHops map[string]map[int]map[int]int) {
	for _, matrix := range matrices {
		for _, node := range matrix.Data {
			setNestedMapValue(discountHops, matrix.Name, node.MicroCategoryID, node.RegionID, node.Price)
		}
	}
}

func (u updateService) PrepareStorage(ctx context.Context, baseLineMatrixName string, discountMatrixNames []string) (models.PreparedStorage, error) {
	var preparedStorage models.PreparedStorage

	baseLineMatrix, err := u.matrixRepo.GetMatrix(ctx, baseLineMatrixName, -1)
	if err != nil {
		return models.PreparedStorage{}, err
	}

	preparedStorage.BaseLineMatrix = baseLineMatrix

	for _, discountMatrixName := range discountMatrixNames {
		discountMatrix, err := u.matrixRepo.GetMatrix(ctx, discountMatrixName, -1)
		if err != nil {
			return models.PreparedStorage{}, err
		}

		preparedStorage.DiscountMatrices = append(preparedStorage.DiscountMatrices, discountMatrix)
	}

	categoryData, regionData, err := u.matrixRepo.GetRelationsWithPrice(ctx, baseLineMatrixName)

	categoryJumps := make([]int, categoryData[len(categoryData)-1][1])
	regionJumps := make([]int, regionData[len(regionData)-1][1])

	recursive(1, categoryData, categoryJumps, 1, 0, true)
	recursive(1, regionData, regionJumps, 1, 0, true)

	preparedStorage.MicroCategoryHops = categoryJumps
	preparedStorage.RegionHops = regionJumps

	firstMap := make(map[int]int)
	firstMap[100] = 100
	secondMap := make(map[int]map[int]int)
	secondMap[-1] = firstMap
	newMap := make(map[string]map[int]map[int]int)
	newMap["-1"] = secondMap

	calculateDiscountMatrix(preparedStorage.DiscountMatrices, newMap)

	preparedStorage.DiscountHops = newMap

	return preparedStorage, nil
}
