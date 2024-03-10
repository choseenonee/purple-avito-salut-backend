package repository

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"template/internal/models"
	"template/pkg/utils"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type matrixRepo struct {
	db *sqlx.DB
}

func InitMatrixRepo(db *sqlx.DB) Matrix {
	return matrixRepo{db: db}
}


func (m matrixRepo) GetHistory(ctx context.Context, timeStart time.Time, timeEnd time.Time, matrixType string) error {

	return nil
}

func (m matrixRepo) CreateMatrix(ctx context.Context, matrix models.MatrixBase) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	valueString := make([]string, 0, len(matrix.Data))
	valueArgs := make([]interface{}, 0, len(matrix.Data))
	for i, node := range matrix.Data {
		valueString = append(valueString, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, matrix.Name, node.MicroCategoryID, node.RegionID, node.Price)
	}

	createMatrixQuery := fmt.Sprintf("INSERT INTO matrix (name, microcategory_id, region_id, price) VALUES %s", strings.Join(valueString, ","))

	res, err := tx.ExecContext(ctx, createMatrixQuery, valueArgs...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}
	if int(count) != len(matrix.Data) {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	createMetadataQuery := `INSERT INTO matrix_metadata (matrix_name, timestamp, is_baseline, parent_matrix_name)
							VALUES ($1, $2, $3, $4);`

	timestamp := time.Now()
	matrixName := fmt.Sprintf("%s_%d", matrix.Name, timestamp.Unix())

	res, err = tx.ExecContext(ctx, createMetadataQuery, matrixName, timestamp, matrix.IsBaseLine, matrix.ParentName)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}
	count, err = res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}
	if int(count) != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}

func (m matrixRepo) GetHistory(ctx context.Context, timeStart time.Time, timeEnd time.Time, isBaseline null.Bool) ([]models.Matrix, error) {
	var matrixes []models.Matrix

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("name, microcategory_id, region_id, matrix_metadata.timestamp, matrix_metadata.parent_matrix_name").
		From("matrix").
		Join("matrix_metadata ON matrix.name = matrix_metadata.matrix_name").
		Where(sq.And{sq.GtOrEq{"matrix_metadata.timestamp": timeStart}, sq.LtOrEq{"matrix_metadata.timestamp": timeEnd}}).OrderBy(`
			matrix_metadata.matrix_name ASC`)

	if isBaseline.Valid {
		query = query.Where(sq.Eq{"matrix_metadata.is_baseline": isBaseline})
	}

	// Собираем запрос
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return []models.Matrix{}, err
	}

	rows, err := m.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return []models.Matrix{}, err
	}

	var matrix models.Matrix

	for rows.Next() {
		var matrixName string
		var matrixTimeStamp time.Time
		var node models.MatrixNode
		var parentMatrixName null.String
		err = rows.Scan(&matrixName, &node.MicroCategoryID, &node.RegionID, &matrixTimeStamp, &parentMatrixName)
		if err != nil {
			return []models.Matrix{}, err
		}

		switch matrix.Name {
		case "":
			matrix.Name = matrixName
			matrix.TimeStamp = matrixTimeStamp
			matrix.ParentName = parentMatrixName
			matrix.Data = append(matrix.Data, node)
		case matrixName:
			matrix.Data = append(matrix.Data, node)
		default:
			matrixes = append(matrixes, matrix)
			matrix.Name = matrixName
			matrix.TimeStamp = matrixTimeStamp
			matrix.ParentName = parentMatrixName
			matrix.Data = nil
			matrix.Data = append(matrix.Data, node)
		}
	}

	if matrix.Name != "" {
		matrixes = append(matrixes, matrix)
	}

	err = rows.Err()
	if err != nil {
		return []models.Matrix{}, nil
	}

	return matrixes, nil
}
