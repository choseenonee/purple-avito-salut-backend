package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"template/pkg/customerr"
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
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createRegionQuery := `INSERT INTO microcategories (name) VALUES ($1) RETURNING id;`

	err = tx.QueryRowxContext(ctx, createRegionQuery, microcategory.Name).Scan(&createdMicrocategoryID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	createRelationshipQuery := `INSERT INTO relationships_microcategories (parent_id, child_id) VALUES ($1, $2);`

	res, err := tx.ExecContext(ctx, createRelationshipQuery, microcategory.ParentID, createdMicrocategoryID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdMicrocategoryID, nil
}
