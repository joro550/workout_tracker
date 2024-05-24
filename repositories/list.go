package repositories

import "database/sql"

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return ListRepository{db: db}
}
