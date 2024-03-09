package fixtures

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
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

	//for key1, value1 := range regions {
	//	fmt.Println(key1)
	//	for key2, value2 := range value1 {
	//		fmt.Println(key2)
	//		for key3, value3 := range value2 {
	//			fmt.Println(key3)
	//			for _, value := range value3 {
	//				fmt.Println(value)
	//			}
	//		}
	//	}
	//}

	return createdIDs, nil
}
