package list

import "database/sql"

type List struct {
	Name        string
	Description string
	Id          int
}

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return ListRepository{db: db}
}
