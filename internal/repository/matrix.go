package repository

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"template/internal/models"
	"template/pkg/customerr"
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

func (m matrixRepo) CreateMatrix(ctx context.Context, matrix models.MatrixBase) (string, error) {
	tx, err := m.db.Beginx()
	if err != nil {
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	timestamp := time.Now()
	matrixName := fmt.Sprintf("%s_%d", matrix.Name, timestamp.Unix())

	var parentNameExists = false
	parentMatrixExistsQuery := `SELECT 1 FROM matrix WHERE name = $1;`
	rows, err := tx.QueryContext(ctx, parentMatrixExistsQuery, matrix.ParentName)
	for rows.Next() {
		parentNameExists = true
	}

	if !parentNameExists {
		return "", customerr.ParentMatrixDontExist
	}

	valueString := make([]string, 0, len(matrix.Data))
	valueArgs := make([]interface{}, 0, len(matrix.Data))
	for i, node := range matrix.Data {
		valueString = append(valueString, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, matrixName, node.MicroCategoryID, node.RegionID, node.Price)
	}

	createMatrixQuery := fmt.Sprintf("INSERT INTO matrix (name, microcategory_id, region_id, price) VALUES %s", strings.Join(valueString, ","))

	res, err := tx.ExecContext(ctx, createMatrixQuery, valueArgs...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if int(count) != len(matrix.Data) {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	createMetadataQuery := `INSERT INTO matrix_metadata (matrix_name, timestamp, is_baseline, parent_matrix_name)
							VALUES ($1, $2, $3, $4);`

	res, err = tx.ExecContext(ctx, createMetadataQuery, matrixName, timestamp, matrix.IsBaseLine, matrix.ParentName)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err = res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if int(count) != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return matrixName, nil
}

func (m matrixRepo) GetHistory(ctx context.Context, data models.GetHistoryMatrix) ([]models.Matrix, error) {
	var matrixes []models.Matrix

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("name, microcategory_id, region_id, matrix_metadata.timestamp, matrix_metadata.parent_matrix_name").
		From("matrix").
		Join("matrix_metadata ON matrix.name = matrix_metadata.matrix_name").
		Where(sq.And{sq.GtOrEq{"matrix_metadata.timestamp": data.TimeStart}, sq.LtOrEq{"matrix_metadata.timestamp": data.TimeEnd}}).OrderBy(`
			matrix_metadata.matrix_name ASC`)

	if data.IsBaseline.Valid {
		query = query.Where(sq.Eq{"matrix_metadata.is_baseline": data.IsBaseline})
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

func (m matrixRepo) GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifference, error) {
	var difference models.MatrixDifference

	deletedAddedQuery := `SELECT matrix.microcategory_id, matrix.region_id, matrix.price
						  FROM matrix
						  INNER JOIN (
						  	SELECT microcategory_id, region_id
						  	FROM matrix
						  	WHERE name=$1
						  	EXCEPT
						  	SELECT microcategory_id, region_id
						  	FROM matrix
						  	WHERE name=$2
						  ) AS subquery ON matrix.microcategory_id = subquery.microcategory_id AND matrix.region_id = subquery.region_id
						  WHERE matrix.name = $1;`

	updatedQuery := `SELECT matrix.microcategory_id, matrix.region_id, matrix.price, m.price
					 FROM matrix
					 JOIN matrix AS m on m.microcategory_id = matrix.microcategory_id AND m.region_id = matrix.region_id
					 WHERE matrix.name=$1 AND m.name=$2 AND matrix.price <> m.price;`

	deletedRows, err := m.db.QueryxContext(ctx, deletedAddedQuery, matrixName1, matrixName2)
	if err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}
	defer deletedRows.Close()

	for deletedRows.Next() {
		var deletedRow models.MatrixNode

		err := deletedRows.Scan(&deletedRow.MicroCategoryID, &deletedRow.RegionID, &deletedRow.Price)
		if err != nil {
			return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}

		difference.Deleted = append(difference.Deleted, deletedRow)
	}

	if err := deletedRows.Err(); err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	addedRows, err := m.db.QueryxContext(ctx, deletedAddedQuery, matrixName2, matrixName1)
	if err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}
	defer addedRows.Close()

	for addedRows.Next() {
		var addedRow models.MatrixNode

		err := addedRows.Scan(&addedRow.MicroCategoryID, &addedRow.RegionID, &addedRow.Price)
		if err != nil {
			return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}

		difference.Added = append(difference.Added, addedRow)
	}

	if err := addedRows.Err(); err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	updatedRows, err := m.db.QueryxContext(ctx, updatedQuery, matrixName1, matrixName2)
	if err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}
	defer updatedRows.Close()

	for updatedRows.Next() {
		var rowBefore models.MatrixNode
		var rowAfter models.MatrixNode

		err := updatedRows.Scan(&rowBefore.MicroCategoryID, &rowBefore.RegionID, &rowBefore.Price, &rowAfter.Price)
		if err != nil {
			return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
		rowAfter.MicroCategoryID = rowBefore.MicroCategoryID
		rowAfter.RegionID = rowBefore.RegionID

		dif := [2]models.MatrixNode{rowBefore, rowAfter}

		difference.Updated = append(difference.Updated, dif)
	}

	if err := updatedRows.Err(); err != nil {
		return models.MatrixDifference{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	return difference, nil
}
