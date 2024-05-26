package task

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/layouts"
	task_pages "github.com/joro550/workout_tracker/task/pages"
	"github.com/joro550/workout_tracker/users"
)

func RegisterTaskEndpoints(mux *chi.Mux, db *sql.DB) {
	taskRepo := NewTaskRepo(db)

	mux.Group(func(r chi.Router) {
		r.Use(users.AuthCtx)

		r.Route("/list/{listId}/task", func(chi chi.Router) {
			chi.Get("/", taskView(taskRepo))
			chi.Delete("/{taskId}", deleteTask(taskRepo))

			chi.Get("/cards", cards(taskRepo))
			chi.Get("/type", taskType)

			chi.Get("/add", taskAddView)
			chi.Post("/add", taskAddPost(taskRepo))
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

		cardView := task_pages.ListTasks(listId, taskList)
		layout := layouts.Authed(cardView)
		page := layouts.Layout(layout)
		page.Render(r.Context(), w)
	}
}

func cards(repo TaskRepository) func(http.ResponseWriter, *http.Request) {
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
		cards := task_pages.Cards(listId, taskList)
		cards.Render(r.Context(), w)
	}
}

func taskType(w http.ResponseWriter, r *http.Request) {
	typeId, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		log.Println("💥 Could not get typeId", err)
	}

	taskType := Type(typeId)
	if WeightSetsAndReps == taskType {
		view := task_pages.WeightSetsAndRepsTempl()
		view.Render(r.Context(), w)
		return
	} else if TimePaceAndDistance == taskType {

		view := task_pages.PaceAndTimeTempl()
		view.Render(r.Context(), w)
		return
	}
}

func taskAddView(w http.ResponseWriter, r *http.Request) {
	listId, _ := strconv.Atoi(r.PathValue("listId"))

	cardView := task_pages.Add(listId)
	layout := layouts.Authed(cardView)
	page := layouts.Layout(layout)
	page.Render(r.Context(), w)
}

func taskAddPost(repo TaskRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		listId, err := strconv.Atoi(r.PathValue("listId"))
		if err != nil {
			log.Println("💥 Could not get listId", err)
		}

		model, err := task_pages.AddTaskModelFromRequest(r)
		if err != nil {
			log.Println("💥 Could not get form information", err)
			return
		}
		_, err = repo.CreateTask(Task{ListId: listId, UserId: user.Id, Title: model.Title, Value: model.Value})
		if err != nil {
			log.Println("💥 Could not create task", err)
		}
		http.Redirect(w, r, fmt.Sprintf("/list/%v/task", listId), http.StatusSeeOther)
	}
}

func deleteTask(repo TaskRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		listId, err := strconv.Atoi(r.PathValue("listId"))
		if err != nil {
			log.Println("💥 Could not get listId", err)
			http.Redirect(w, r, fmt.Sprintf("/list/%v/task", listId), http.StatusSeeOther)
			return
		}

		taskId, err := strconv.Atoi(r.PathValue("taskId"))
		if err != nil {
			log.Println("💥 Could not get listId", err)
			http.Redirect(w, r, fmt.Sprintf("/list/%v/task", listId), http.StatusSeeOther)
			return
		}

		_, err = repo.DeleteTask(taskId, listId, user.Id)
		if err != nil {
			log.Println("💥 Could not delete task", err)
			http.Redirect(w, r, fmt.Sprintf("/list/%v/task", listId), http.StatusSeeOther)
			return
		}

		w.Header().Set("HX-Trigger", "task-updated")
		w.WriteHeader(http.StatusOK)
	}
}
