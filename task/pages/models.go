package task_pages

import "time"

type TaskModel struct {
	Title string
	Date  time.Time
	Value string
	Id    int
}
