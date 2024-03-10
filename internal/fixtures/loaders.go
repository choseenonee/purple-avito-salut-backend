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
	Regions         = "../internal/fixtures/regions.json"
	UsersSegments   = "../internal/fixtures/users-segments.json"
	Users           = "../internal/fixtures/users.json"
	Microcategories = "../internal/fixtures/microcategories.json"
)

func Loader(filepath ...string) error {
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

		switch path {
		case Regions:
			regionRepo := repository.InitRegionsRepo(db)

			var rawRegions map[string]map[string]map[string][]string
			if err := json.Unmarshal(bytes, &rawRegions); err != nil {
				return err
			}

			_, err := loadRegions(regionRepo, rawRegions)
			if err != nil {
				return err
			}
			//fmt.Println(createdIDs)
		case UsersSegments:
			usersSegmentsRepo := repository.InitSegmentsRepo(db)

			var arrayUsersSegments []parsingUsersSegments
			if err := json.Unmarshal(bytes, &arrayUsersSegments); err != nil {
				return err
			}

			_, err := loadUsersSegments(usersSegmentsRepo, arrayUsersSegments)
			if err != nil {
				return err
			}
			//fmt.Println(createdIDs)
		case Users:
			userRepo := repository.InitUserRepo(db)

			var rawUsers []models.UserBase
			if err := json.Unmarshal(bytes, &rawUsers); err != nil {
				return err
			}

			_, err := loadUsers(userRepo, rawUsers)
			if err != nil {
				return err
			}
			//fmt.Println(createdIDs)
		case Microcategories:
			microcategoriesRepo := repository.InitMicrocategoryRepo(db)

			var rawMicrocategories map[string]map[string]map[string]map[string][]string
			if err := json.Unmarshal(bytes, &rawMicrocategories); err != nil {
				return err
			}

			_, err := loadMicrocategories(microcategoriesRepo, rawMicrocategories)
			if err != nil {
				return err
			}
			//fmt.Println(createdIDs)
		default:
			return fmt.Errorf("error parsing path")
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

type parsingUsersSegments struct {
	UserID     int   `json:"user_id"`
	SegmentIDs []int `json:"segment_id"`
}

func loadUsersSegments(usersSegmentsRepo repository.UsersSegments, usersSegments []parsingUsersSegments) ([]int, error) {
	var createdIDs []int

	for _, userSegmentNode := range usersSegments {
		for _, userSegmentRaw := range userSegmentNode.SegmentIDs {
			userSegment := models.UserSegmentBase{
				UserID:    userSegmentNode.UserID,
				SegmentID: userSegmentRaw,
			}
			id, err := usersSegmentsRepo.Create(context.Background(), userSegment)
			if err != nil {
				return []int{}, err
			}
			createdIDs = append(createdIDs, id)
		}
	}

	return createdIDs, nil
}

func loadUsers(userRepo repository.Users, users []models.UserBase) ([]int, error) {
	var createdIDs []int

	for _, user := range users {
		id, err := userRepo.Create(context.Background(), user)
		if err != nil {
			return []int{}, err
		}

		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

func loadMicrocategories(microcategoryRepo repository.Microcategories, microcategories map[string]map[string]map[string]map[string][]string) ([]int, error) {
	var createdIDs []int

	for key1, value1 := range microcategories {
		data := models.MicrocategoryBase{ParentID: null.NewInt(0, false), Name: key1}
		id1, err := microcategoryRepo.Create(context.Background(), data)
		if err != nil {
			return []int{}, err
		}

		for key2, value2 := range value1 {
			data := models.MicrocategoryBase{ParentID: null.NewInt(int64(id1), true), Name: key2}
			id2, err := microcategoryRepo.Create(context.Background(), data)
			if err != nil {
				return []int{}, err
			}

			for key3, value3 := range value2 {
				data := models.MicrocategoryBase{ParentID: null.NewInt(int64(id2), true), Name: key3}
				id3, err := microcategoryRepo.Create(context.Background(), data)
				if err != nil {
					return []int{}, err
				}

				for key4, value4 := range value3 {
					data := models.MicrocategoryBase{ParentID: null.NewInt(int64(id3), true), Name: key4}
					id4, err := microcategoryRepo.Create(context.Background(), data)
					if err != nil {
						return []int{}, err
					}

					for _, value := range value4 {
						data := models.MicrocategoryBase{ParentID: null.NewInt(int64(id4), true), Name: value}
						id5, err := microcategoryRepo.Create(context.Background(), data)
						if err != nil {
							return []int{}, err
						}

						createdIDs = append(createdIDs, id5)
					}

					createdIDs = append(createdIDs, id3)
				}

				createdIDs = append(createdIDs, id2)
			}

			createdIDs = append(createdIDs, id1)
		}
	}

	return createdIDs, nil
}
