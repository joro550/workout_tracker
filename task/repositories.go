package task

import (
	"database/sql"
	"log"
	"time"
)

type Type int

const (
	WeightSetsAndReps Type = iota
	TimePaceAndDistance
)

type Task struct {
	Title string
	Date  time.Time
	Value string
	Type  Type
	Id    int
}

type WeightTask struct{}

type TimeTask struct{}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepository {
	return TaskRepository{db: db}
}

func (tr TaskRepository) GetAllTasks(userId, listId int) ([]Task, error) {
	rows, err := tr.db.Query("select id, title, date, value, type from task where userid = ? and id = ? ",
		userId, listId)
	if err != nil {
		log.Println("🤔 [GetAllTasks] query faile", err)
		return []Task{}, nil
	}

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Title, &task.Date, &task.Value, &task.Type)
		if err != nil {
			log.Println("💥 [GetAllTasks] scan failed", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
