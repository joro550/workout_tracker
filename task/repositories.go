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
	Title  string
	Date   time.Time
	Value  string
	Type   Type
	Id     int
	ListId int
	UserId int
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
	rows, err := tr.db.Query("select id, title, date, value, type, listid, userid from task where userid = ? and listid = ? ",
		userId, listId)
	if err != nil {
		log.Println("🤔 [GetAllTasks] query faile", err)
		return []Task{}, nil
	}

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Title, &task.Date, &task.Value, &task.Type, &task.ListId, &task.UserId)
		if err != nil {
			log.Println("💥 [GetAllTasks] scan failed", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (tr TaskRepository) CreateTask(task Task) (int, error) {
	rows := tr.db.QueryRow(
		`insert into task (title, date, value, type, listId, userId)
    values (?,?,?,?,?,?)
    returning id`,
		task.Title, time.Now(), task.Value, task.Type, task.ListId, task.UserId)

	var id int
	err := rows.Scan(&id)
	return id, err
}

func (tr TaskRepository) DeleteTask(id, listId, userId int) (bool, error) {
	_, err := tr.db.Exec("delete from task where id = ? and listid = ? and userid = ?", id, listId, userId)
	if err != nil {
		log.Println("🤔 [DeleteTask] query failed to execure", err)
		return false, err
	}
	return true, nil
}
