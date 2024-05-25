package task

import "database/sql"

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepository {
	return TaskRepository{db: db}
}
