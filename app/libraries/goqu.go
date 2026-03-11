package libraries

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

type goquLibrary struct {
	db *goqu.Database
}

type GoquLibrary interface {
	DB() *goqu.Database
}

func NewGoquLibrary(db *sql.DB) *goquLibrary {
	return &goquLibrary{
		db: goqu.New("mysql", db),
	}
}

func (g *goquLibrary) DB() *goqu.Database {
	return g.db
}