package repository

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"template/internal/models"
)

type nodesRepo struct {
	db *sqlx.DB
}

func InitNodesRepo(db *sqlx.DB) Nodes {
	return nodesRepo{db: db}
}

func (n nodesRepo) GetMicrocategoryTree() ([]models.NodeRaw, error) {
	var nodes []models.NodeRaw

	//gets PARENT price
	getMicrocategoryTreeQuery := `SELECT parent_id, child_id, matrix.price FROM relationships_microcategories
									LEFT JOIN matrix ON parent_id = matrix.microcategory_id`

	rows, err := n.db.Query(getMicrocategoryTreeQuery)
	if err != nil {
		return nil, nil
	}

	var children []int
	var node models.NodeRaw

	for rows.Next() {
		var nodeID int
		var childID int
		var price null.Int

		err = rows.Scan(&nodeID, &childID, &price)
		if err != nil {
			return nil, nil
		}

		switch node.ID {
		case 0:
			node.ID = nodeID
			if price.Valid {
				node.HasPrice = true
			}
		case nodeID:
			children = append(children, childID)
		default:
			node.ChildrenIDs = children
			//copy(node.ChildrenIDs, children)
			nodes = append(nodes, node)
			node.ID = 0
			children = nil
		}

	}

	nodes = append(nodes, node)

	return nodes, nil
}
