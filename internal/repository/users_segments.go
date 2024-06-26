package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"template/pkg/customerr"
)

type usersSegmentsRepo struct {
	db *sqlx.DB
}

func InitSegmentsRepo(db *sqlx.DB) UsersSegments {
	return usersSegmentsRepo{db: db}
}

func (r usersSegmentsRepo) Create(ctx context.Context, userSegment models.UserSegmentBase) (int, error) {
	var createdUserSegmentID int

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createRegionQuery := `INSERT INTO users_segments (user_id, segment_id) VALUES ($1, $2) RETURNING id;`

	err = tx.QueryRowxContext(ctx, createRegionQuery, userSegment.UserID, userSegment.SegmentID).Scan(&createdUserSegmentID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdUserSegmentID, nil
}
