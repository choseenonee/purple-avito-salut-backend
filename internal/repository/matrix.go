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
	db        *sqlx.DB
	MaxOnPage int
}

func InitMatrixRepo(db *sqlx.DB, MaxOnPage int) Matrix {
	return matrixRepo{MaxOnPage: MaxOnPage, db: db}
}

func (m matrixRepo) CreateMatrixWithoutParent(ctx context.Context, matrix models.MatrixBase) (string, error) {
	tx, err := m.db.Beginx()
	if err != nil {
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	timestamp := time.Now()
	matrixName := fmt.Sprintf("%s_%d", matrix.Name, timestamp.Unix())

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

func (m matrixRepo) CreateMatrix(ctx context.Context, matrix models.MatrixDifferenceRequest) (string, error) {
	tx, err := m.db.Beginx()
	if err != nil {
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	timestamp := time.Now()
	matrixName := fmt.Sprintf("%s_%d", matrix.NewName, timestamp.Unix())

	var parentNameExists = false

	parentMatrixExistsQuery := `SELECT 1 FROM matrix WHERE name = $1;`
	rows, err := tx.QueryContext(ctx, parentMatrixExistsQuery, matrix.ParentName)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.QueryRrr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}

	defer rows.Close()

	for rows.Next() {
		parentNameExists = true
	}
	if !parentNameExists {
		return "", customerr.ParentMatrixDontExist
	}

	insertCopyQuery := `INSERT INTO matrix (name, microcategory_id, region_id, price)
						SELECT $1, microcategory_id, region_id, price FROM matrix WHERE name = $2;`

	_, err = tx.ExecContext(ctx, insertCopyQuery, matrixName, matrix.ParentName)
	if err != nil {
		return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	for _, row := range matrix.Updated {
		updateQuery := `UPDATE matrix SET price = $1
              			WHERE name = $2 AND microcategory_id = $3 AND region_id = $4;`
		res, err := tx.ExecContext(ctx, updateQuery, row.Price, matrixName, row.MicroCategoryID, row.RegionID)
		if err != nil {
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
		if int(count) != 1 {
			if rbErr := tx.Rollback(); rbErr != nil {
				return "", customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}
			return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
		}
	}

	for _, row := range matrix.Deleted {
		deleteQuery := `DELETE FROM matrix
						WHERE name = $1 AND microcategory_id = $2 AND region_id = $3;`
		_, err = tx.ExecContext(ctx, deleteQuery, matrixName, row.MicroCategoryID, row.RegionID)
		if err != nil {
			return "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
		}
	}

	valueString := make([]string, 0, len(matrix.Added))
	valueArgs := make([]interface{}, 0, len(matrix.Added))
	for i, node := range matrix.Added {
		valueString = append(valueString, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, matrixName, node.MicroCategoryID, node.RegionID, node.Price)
	}

	createRowsQuery := fmt.Sprintf("INSERT INTO matrix (name, microcategory_id, region_id, price) VALUES %s", strings.Join(valueString, ","))
	res, err := tx.ExecContext(ctx, createRowsQuery, valueArgs...)
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
	if int(count) != len(matrix.Added) {
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

func (m matrixRepo) GetMatrixPages(ctx context.Context, matrixName string) (int, error) {
	var count int

	getMatrixPagesQuery := "SELECT COUNT(*) FROM matrix WHERE name=$1"

	err := m.db.QueryRowxContext(ctx, getMatrixPagesQuery, matrixName).Scan(&count)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return (count + m.MaxOnPage - 1) / m.MaxOnPage, nil
}

func (m matrixRepo) GetHistory(ctx context.Context, data models.GetHistoryMatrix) ([]models.ResponseHistoryMatrix, error) {
	var matrixes []models.ResponseHistoryMatrix

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("DISTINCT name, matrix_metadata.timestamp, matrix_metadata.parent_matrix_name").
		From("matrix").
		Join("matrix_metadata ON matrix.name = matrix_metadata.matrix_name").
		Where(sq.And{sq.GtOrEq{"matrix_metadata.timestamp": data.TimeStart}, sq.LtOrEq{"matrix_metadata.timestamp": data.TimeEnd}})

	if data.IsBaseline.Valid {
		query = query.Where(sq.Eq{"matrix_metadata.is_baseline": data.IsBaseline})
	}

	// Собираем запрос
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return []models.ResponseHistoryMatrix{}, err
	}

	rows, err := m.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return []models.ResponseHistoryMatrix{}, err
	}

	defer rows.Close()

	var matrix models.ResponseHistoryMatrix

	for rows.Next() {
		err = rows.Scan(&matrix.Name, &matrix.TimeStamp, &matrix.ParentName)
		if err != nil {
			return []models.ResponseHistoryMatrix{}, err
		}

		matrixes = append(matrixes, matrix)
	}

	err = rows.Err()
	if err != nil {
		return []models.ResponseHistoryMatrix{}, nil
	}

	return matrixes, nil
}

func (m matrixRepo) GetPriceTendency(ctx context.Context, data models.GetTendencyNode) ([]models.ResponseTendencyNode, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("DISTINCT matrix_metadata.timestamp, price").
		From("matrix").
		Join("matrix_metadata ON matrix.name = matrix_metadata.matrix_name").
		Where(sq.And{sq.GtOrEq{"matrix_metadata.timestamp": data.TimeStart},
			sq.LtOrEq{"matrix_metadata.timestamp": data.TimeEnd}, sq.Eq{"microcategory_id": data.MicrocategoryID},
			sq.Eq{"region_id": data.RegionID}, sq.Eq{"matrix_metadata.is_baseline": true}})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return []models.ResponseTendencyNode{}, err
	}

	rows, err := m.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return []models.ResponseTendencyNode{}, err
	}

	defer rows.Close()

	var responses []models.ResponseTendencyNode

	for rows.Next() {
		var response models.ResponseTendencyNode
		err = rows.Scan(&response.TimeStamp, &response.Price)
		if err != nil {
			return []models.ResponseTendencyNode{}, err
		}

		responses = append(responses, response)
	}

	err = rows.Err()
	if err != nil {
		return []models.ResponseTendencyNode{}, nil
	}

	oneBeforeTimeStartQuery := `SELECT DISTINCT matrix_metadata.timestamp, price 
								FROM matrix 
    							JOIN matrix_metadata ON matrix.name = matrix_metadata.matrix_name 
								WHERE matrix_metadata.timestamp <= $1 AND matrix_metadata.is_baseline = true 
								  AND microcategory_id = $2 AND region_id = $3
								ORDER BY matrix_metadata.timestamp DESC`

	rows.Close()

	rows, err = m.db.QueryxContext(ctx, oneBeforeTimeStartQuery, data.TimeStart, data.MicrocategoryID, data.RegionID)
	if err != nil {
		return []models.ResponseTendencyNode{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var response models.ResponseTendencyNode
		err = rows.Scan(&response.TimeStamp, &response.Price)
		if err != nil {
			return []models.ResponseTendencyNode{}, err
		}

		responses = append(responses, response)
	}

	err = rows.Err()
	if err != nil {
		return []models.ResponseTendencyNode{}, nil
	}

	return responses, nil
}

func (m matrixRepo) GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifferenceResponse, error) {
	var difference models.MatrixDifferenceResponse

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
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer deletedRows.Close()

	for deletedRows.Next() {
		var deletedRow models.MatrixNode

		err := deletedRows.Scan(&deletedRow.MicroCategoryID, &deletedRow.RegionID, &deletedRow.Price)
		if err != nil {
			return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		difference.Deleted = append(difference.Deleted, deletedRow)
	}

	if err := deletedRows.Err(); err != nil {
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	addedRows, err := m.db.QueryxContext(ctx, deletedAddedQuery, matrixName2, matrixName1)
	if err != nil {
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer addedRows.Close()

	for addedRows.Next() {
		var addedRow models.MatrixNode

		err := addedRows.Scan(&addedRow.MicroCategoryID, &addedRow.RegionID, &addedRow.Price)
		if err != nil {
			return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		difference.Added = append(difference.Added, addedRow)
	}

	if err := addedRows.Err(); err != nil {
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	updatedRows, err := m.db.QueryxContext(ctx, updatedQuery, matrixName1, matrixName2)
	if err != nil {
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer updatedRows.Close()

	for updatedRows.Next() {
		var rowBefore models.MatrixNode
		var rowAfter models.MatrixNode

		err := updatedRows.Scan(&rowBefore.MicroCategoryID, &rowBefore.RegionID, &rowBefore.Price, &rowAfter.Price)
		if err != nil {
			return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		rowAfter.MicroCategoryID = rowBefore.MicroCategoryID
		rowAfter.RegionID = rowBefore.RegionID

		dif := [2]models.MatrixNode{rowBefore, rowAfter}

		difference.Updated = append(difference.Updated, dif)
	}

	if err := updatedRows.Err(); err != nil {
		return models.MatrixDifferenceResponse{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return difference, nil
}

func (m matrixRepo) GetMatrix(ctx context.Context, matrixName string, mc, rg null.Int, page int) (models.Matrix, error) {
	var matrix models.Matrix
	var rows *sqlx.Rows
	var err error

	if page == -1 {
		selectQuery := `SELECT name, microcategory_id, region_id, price, mm.timestamp, mm.is_baseline, mm.parent_matrix_name FROM matrix
					JOIN matrix_metadata mm ON matrix.name = mm.matrix_name
                    WHERE name = $1
                    ORDER BY (microcategory_id, region_id) DESC`

		rows, err = m.db.QueryxContext(ctx, selectQuery, matrixName)

		if err != nil {
			return models.Matrix{}, err
		}
	} else {
		selectQuery := `SELECT name, microcategory_id, region_id, price, mm.timestamp, mm.is_baseline, mm.parent_matrix_name FROM matrix
						JOIN matrix_metadata mm ON matrix.name = mm.matrix_name`

		queryParams := []interface{}{matrixName}
		whereClauses := []string{"name = $1"}
		if mc.Valid {
			whereClauses = append(whereClauses, fmt.Sprintf("microcategory_id = $%d", len(queryParams)+1))
			queryParams = append(queryParams, mc.Int64) // Добавляем значение mc в параметры запроса
		}
		if rg.Valid {
			whereClauses = append(whereClauses, fmt.Sprintf("region_id = $%d", len(queryParams)+1))
			queryParams = append(queryParams, rg.Int64) // Добавляем значение rg в параметры запроса
		}

		selectQuery += " WHERE " + strings.Join(whereClauses, " AND ")

		offset := (page - 1) * m.MaxOnPage
		selectQuery += fmt.Sprintf(" ORDER BY (microcategory_id, region_id) DESC OFFSET $%d LIMIT $%d", len(queryParams)+1, len(queryParams)+2)
		queryParams = append(queryParams, offset, m.MaxOnPage)
		rows, err = m.db.QueryxContext(ctx, selectQuery, queryParams...)

		if err != nil {
			return models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
		}
	}

	defer rows.Close()

	for rows.Next() {
		var row models.MatrixNode

		err := rows.Scan(&matrix.Name, &row.MicroCategoryID, &row.RegionID, &row.Price, &matrix.TimeStamp, &matrix.IsBaseLine, &matrix.ParentName)
		if err != nil {
			return models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		matrix.Data = append(matrix.Data, row)
	}

	if err := rows.Err(); err != nil {
		return models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	matrix.Name = matrixName

	return matrix, nil
}

func (m matrixRepo) GetMatricesByDuration(ctx context.Context, timeStart, timeEnd time.Time) ([]models.Matrix, error) {
	var matrices []models.Matrix

	selectQuery := `SELECT matrix_name, timestamp, is_baseline, parent_matrix_name FROM matrix_metadata
					WHERE timestamp BETWEEN $1 AND $2;`

	rows, err := m.db.QueryxContext(ctx, selectQuery, timeStart, timeEnd)
	if err != nil {
		return []models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var matrix models.Matrix

		err := rows.Scan(&matrix.Name, &matrix.TimeStamp, &matrix.IsBaseLine, &matrix.ParentName)
		if err != nil {
			return []models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		matrices = append(matrices, matrix)
	}

	if err := rows.Err(); err != nil {
		return []models.Matrix{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return matrices, nil
}

func (m matrixRepo) GetRelationsWithPrice(ctx context.Context, matrixName string) ([][4]int, [][4]int, error) {
	var categoryData [][4]int
	var regionData [][4]int

	categoryQuery := `SELECT
		rr.parent_id,
		rr.child_id,
		matrix_parent.price AS parent_price,
		matrix_child.price AS child_price
	FROM relationships_microcategories rr
	LEFT JOIN matrix AS matrix_parent ON rr.parent_id = matrix_parent.microcategory_id AND matrix_parent.name = $1
	LEFT JOIN matrix AS matrix_child ON rr.child_id = matrix_child.microcategory_id AND matrix_child.name = $1
	ORDER BY rr.parent_id;`

	rows, err := m.db.QueryContext(ctx, categoryQuery, matrixName)
	if err != nil {
		return [][4]int{}, [][4]int{}, nil
	}

	defer rows.Close()

	for rows.Next() {
		var parentID int
		var childID int
		var parentPrice null.Int
		var childPrice null.Int

		err = rows.Scan(&parentID, &childID, &parentPrice, &childPrice)
		if err != nil {
			return [][4]int{}, [][4]int{}, nil
		}

		categoryData = append(categoryData, [4]int{parentID, childID, int(parentPrice.Int64), int(childPrice.Int64)})
	}

	if rows.Err() != nil {
		return [][4]int{}, [][4]int{}, nil
	}

	regionQuery := `SELECT
	rr.parent_id,
		rr.child_id,
		matrix_parent.price AS parent_price,
		matrix_child.price AS child_price
	FROM
	relationships_regions rr
	LEFT JOIN
	matrix AS matrix_parent ON rr.parent_id = matrix_parent.region_id AND matrix_parent.name = $1
	LEFT JOIN
	matrix AS matrix_child ON rr.child_id = matrix_child.region_id AND matrix_child.name = $1
	ORDER BY
	rr.parent_id;`

	rows, err = m.db.QueryContext(ctx, regionQuery, matrixName)
	if err != nil {
		return [][4]int{}, [][4]int{}, nil
	}

	defer rows.Close()

	for rows.Next() {
		var parentID int
		var childID int
		var parentPrice null.Int
		var childPrice null.Int

		err = rows.Scan(&parentID, &childID, &parentPrice, &childPrice)
		if err != nil {
			return [][4]int{}, [][4]int{}, nil
		}

		regionData = append(regionData, [4]int{parentID, childID, int(parentPrice.Int64), int(childPrice.Int64)})
	}

	if rows.Err() != nil {
		return [][4]int{}, [][4]int{}, nil
	}

	return categoryData, regionData, nil
}
