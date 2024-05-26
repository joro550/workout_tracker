package task_pages

import (
	"net/http"
	"time"
)

type TaskModel struct {
	Title string
	Date  time.Time
	Value string
	Id    int
}

type AddTaskModel struct {
	Title string
	Value string
}

func AddTaskModelFromRequest(r *http.Request) AddTaskModel {
	return AddTaskModel{
		Title: r.FormValue("title"),
	}
}
