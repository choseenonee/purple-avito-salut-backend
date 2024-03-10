package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"template/pkg/utils"
)

type microcategoryRepo struct {
	db *sqlx.DB
}

func InitMicrocategoryRepo(db *sqlx.DB) Microcategories {
	return microcategoryRepo{db: db}
}

func (m microcategoryRepo) Create(ctx context.Context, microcategory models.MicrocategoryBase) (int, error) {
	var createdMicrocategoryID int

	tx, err := m.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	createRegionQuery := `INSERT INTO microcategories (name) VALUES ($1) RETURNING id;`

	err = tx.QueryRowxContext(ctx, createRegionQuery, microcategory.Name).Scan(&createdMicrocategoryID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	createRelationshipQuery := `INSERT INTO relationships_microcategories (parent_id, child_id) VALUES ($1, $2);`

	res, err := tx.ExecContext(ctx, createRelationshipQuery, microcategory.ParentID, createdMicrocategoryID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdMicrocategoryID, nil
}
