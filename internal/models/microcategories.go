package models

import "github.com/guregu/null"

type MicrocategoryBase struct {
	ParentID null.Int
	Name     string
}
