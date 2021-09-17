package postgres

import "database/sql"

type Store struct {
	*sql.DB
}
