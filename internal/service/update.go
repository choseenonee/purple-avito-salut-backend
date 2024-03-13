package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"template/internal/models"
	"template/internal/repository"
)

type updateService struct {
	matrixRepo     repository.Matrix
	currentStorage *models.PreparedStorageSend
}

func InitUpdateService(matrixRepo repository.Matrix) Update {
	return &updateService{matrixRepo: matrixRepo}
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

var segmentMatrices = map[int]string{
	100: "discount_0",
	200: "discount_1",
}

func (u *updateService) PrepareStorage(ctx context.Context, baseLineMatrixName string, discountMatrixNames []string) (models.PreparedStorage, error) {
	var preparedStorage models.PreparedStorage

	preparedStorage.SegmentDiscount = segmentMatrices

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

func (u *updateService) SendUpdatedStorage(url string, storage models.PreparedStorageSend) error {
	jsonData, err := json.Marshal(storage)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/update_next_storage", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error response status: %v", resp.StatusCode)
	}

	u.currentStorage = &models.PreparedStorageSend{
		StorageBase: models.StorageBase{
			BaseLineMatrixName:  storage.BaseLineMatrixName,
			DiscountMatrixNames: storage.DiscountMatrixNames,
		},
		MicroCategoryHops: storage.MicroCategoryHops,
		RegionHops:        storage.RegionHops,
		DiscountHops:      storage.DiscountHops,
		SegmentDiscount:   storage.SegmentDiscount,
	}

	return nil
}

func (u *updateService) SwitchStorage(url string) error {
	url = fmt.Sprintf("%v/update_current_storage", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error response status: %v", resp.StatusCode)
	}
	return nil
}

func (u *updateService) GetCurrentStorage() models.PreparedStorageSend {
	if u.currentStorage == nil {
		return models.PreparedStorageSend{}
	}
	return *u.currentStorage
}
