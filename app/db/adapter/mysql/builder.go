package mysql

import (
	"database/sql"
	"gorm.io/gorm"
	"keeper/app/db/standard/modules"
)

type Query struct {
	Rows    *sql.Rows
	Columns []*modules.Column `json:"columns"`
}

func getSqlQuery(db *gorm.DB, sql string) (*Query, error) {
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return nil, nil
	}
	return &Query{Rows: rows, Columns: getSqlColumns(rows)}, nil
}
