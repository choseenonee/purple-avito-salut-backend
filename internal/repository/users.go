package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"template/pkg/utils"
)

type userRepo struct {
	db *sqlx.DB
}

func InitUserRepo(db *sqlx.DB) Users {
	return userRepo{db: db}
}

func (u userRepo) Create(ctx context.Context, user models.UserBase) (int, error) {
	var createdUserID int

	tx, err := u.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	createUserQuery := `INSERT INTO users (name, region_id) VALUES ($1, $2) RETURNING id;`

	err = tx.QueryRowxContext(ctx, createUserQuery, user.Name, user.RegionID).Scan(&createdUserID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdUserID, nil
}
