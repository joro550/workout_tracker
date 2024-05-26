package task

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joro550/workout_tracker/layouts"
	task_pages "github.com/joro550/workout_tracker/task/pages"
	"github.com/joro550/workout_tracker/users"
)

func RegisterTaskEndpoints(mux *chi.Mux, db *sql.DB) {
	taskRepo := NewTaskRepo(db)
	mux.Group(func(r chi.Router) {
		r.Use(users.AuthCtx)

		r.Route("/list/{listId}/task", func(chi chi.Router) {
			r.Get("/", taskView(taskRepo))
		})
	})
}

func taskView(repo TaskRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		listId, err := strconv.Atoi(r.PathValue("listId"))
		if err != nil {
			log.Println("💥 Could not get listId", err)
		}

		tasks, err := repo.GetAllTasks(user.Id, listId)
		if err != nil {
			log.Println("💥 Could not get task", err)
		}

		var taskList []task_pages.TaskModel
		for _, task := range tasks {
			taskList = append(taskList, task_pages.TaskModel{
				Id:    task.Id,
				Value: task.Value,
				Title: task.Title,
				Date:  task.Date,
			})
		}

		cardView := task_pages.UpdateableCards(taskList)
		layout := layouts.Authed(cardView)
		page := layouts.Layout(layout)
		page.Render(r.Context(), w)
	}
}
