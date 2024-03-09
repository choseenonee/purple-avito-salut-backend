package fixtures

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"os"
	"template/internal/models"
	"template/internal/repository"
	"template/pkg/config"
	"template/pkg/database"
)

const (
	Regions = "../internal/fixtures/regions.json"
)

func LoadRegions(filepath ...string) error {
	var db *sqlx.DB

	if len(filepath) > 0 {
		config.InitConfig()
		db = database.GetDB()
	}

	for _, path := range filepath {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var rawRegions map[string]map[string]map[string][]string
		if err := json.Unmarshal(bytes, &rawRegions); err != nil {
			return err
		}

		switch path {
		case Regions:
			regionRepo := repository.InitRegionsRepo(db)

			createdIDs, err := loadRegions(regionRepo, rawRegions)
			if err != nil {
				return err
			}
			fmt.Println(createdIDs)
		default:
			return fmt.Errorf("иди нахуй")
		}

	}

	return nil
}

func loadRegions(regionRepo repository.Regions, regions map[string]map[string]map[string][]string) ([]int, error) {
	var createdIDs []int

	for key1, value1 := range regions {
		data := models.RegionBase{ParentID: null.NewInt(0, false), Name: key1}
		id1, err := regionRepo.Create(context.Background(), data)
		if err != nil {
			return []int{}, err
		}

		for key2, value2 := range value1 {
			data := models.RegionBase{ParentID: null.NewInt(int64(id1), true), Name: key2}
			id2, err := regionRepo.Create(context.Background(), data)
			if err != nil {
				return []int{}, err
			}

			for key3, value3 := range value2 {
				data := models.RegionBase{ParentID: null.NewInt(int64(id2), true), Name: key3}
				id3, err := regionRepo.Create(context.Background(), data)
				if err != nil {
					return []int{}, err
				}

				for _, value := range value3 {
					data := models.RegionBase{ParentID: null.NewInt(int64(id3), true), Name: value}
					id4, err := regionRepo.Create(context.Background(), data)
					if err != nil {
						return []int{}, err
					}

					createdIDs = append(createdIDs, id4)
				}

				createdIDs = append(createdIDs, id3)
			}

			createdIDs = append(createdIDs, id2)
		}

		createdIDs = append(createdIDs, id1)
	}

	return createdIDs, nil
}
