package repository

import (
	"context"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type matrixRepo struct {
	db *sqlx.DB
}

func InitmatrixRepo(db *sqlx.DB) Matrix {
	return matrixRepo{db: db}
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
