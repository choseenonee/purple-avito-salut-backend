package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
	"time"
)

type matrixRepo struct {
	db *sqlx.DB
}

func InitmatrixRepo(db *sqlx.DB) Matrix {
	return matrixRepo{db: db}
}

func (m matrixRepo) GetHistory(ctx context.Context, timeStart time.Time, timeEnd time.Time, matrixType string) {
	var matrixes []models.Matrix

	query := `SELECT matrix_name, microcategory_id, region_id, matrix_metadata.timestamp 
				FROM matrix
				JOIN matrix_metadata ON matrix.name = matrix_metadata.name AND matrix_metadata.timestamp BETWEEN ($1, $2) AND matrix_metadata.type = `
}
