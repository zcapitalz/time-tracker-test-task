package storages

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func (s *SQLStorage) Init(db *sqlx.DB) {
	if s == nil {
		panic("can't init nil storage")
	}

	s.db = db
	s.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
