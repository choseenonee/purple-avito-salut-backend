package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"template/pkg/customerr"
)

type Repository interface {
	GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error)
	//GetRegionPath(regionID int) ([]int, error)
	//// string - matrix name (from where we got this)
	//GetPriceFromBaseLine(microcategoryID int, regionID int) (int, string, error)
	// TODO: discount matrix add
}

type repostitoryStruct struct {
	db *sqlx.DB
}

func InitRepository(db *sqlx.DB) Repository {
	return repostitoryStruct{db: db}
}

func (r repostitoryStruct) GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error) {
	path := make([]int, 0, 10)

	selectQuery := `WITH RECURSIVE path AS (
    SELECT
        child_id,
        parent_id,
        ARRAY[child_id] AS path_array -- Массив для хранения пути
    FROM relationships_microcategories
    WHERE child_id = $1 -- Замените :your_child_id на интересующий вас child_id
    UNION ALL
    SELECT
        rl.child_id,
        rl.parent_id,
        p.path_array || rl.child_id -- Добавляем child_id к пути
    FROM
        relationships_microcategories rl
            JOIN path p ON p.parent_id = rl.child_id
)
SELECT
    path_array
FROM
    path
WHERE
    parent_id = 1;`

	rows, err := r.db.QueryxContext(ctx, selectQuery, microCategoryID)
	if err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer rows.Close()

	var rawPath []sql.NullInt32
	for rows.Next() {
		err := rows.Scan(pq.Array(&rawPath))
		if err != nil {
			return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
	}

	for _, i := range rawPath {
		if i.Valid {
			path = append(path, int(i.Int32))
		}
	}

	if err := rows.Err(); err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return path, nil
}

func (r repostitoryStruct) GetRegionPath(ctx context.Context, microCategoryID int) ([]int, error) {
	path := make([]int, 0, 10)

	selectQuery := `WITH RECURSIVE path AS (
    SELECT
        child_id,
        parent_id,
        ARRAY[child_id] AS path_array -- Массив для хранения пути
    FROM relationships_microcategories
    WHERE child_id = $1 -- Замените :your_child_id на интересующий вас child_id
    UNION ALL
    SELECT
        rl.child_id,
        rl.parent_id,
        p.path_array || rl.child_id -- Добавляем child_id к пути
    FROM
        relationships_microcategories rl
            JOIN path p ON p.parent_id = rl.child_id
)
SELECT
    path_array
FROM
    path
WHERE
    parent_id = 1;`

	rows, err := r.db.QueryxContext(ctx, selectQuery, microCategoryID)
	if err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryRrr, Err: err})
	}
	defer rows.Close()

	var rawPath []sql.NullInt32
	for rows.Next() {
		err := rows.Scan(pq.Array(&rawPath))
		if err != nil {
			return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
	}

	for _, i := range rawPath {
		if i.Valid {
			path = append(path, int(i.Int32))
		}
	}

	if err := rows.Err(); err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return path, nil
}
