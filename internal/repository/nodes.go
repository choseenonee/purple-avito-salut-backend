package repository

import (
	"github.com/jmoiron/sqlx"
	"template/internal/models"
)

type nodesRepo struct {
	db *sqlx.DB
}

func InitNodesRepo(db *sqlx.DB) Nodes {
	return nodesRepo{db: db}
}

func (n nodesRepo) GetMicrocategoryTree() ([]models.Node, error) {
	var nodes []models.Node

	getMicrocategoryTreeQuery := ``

	return nodes, nil
}
