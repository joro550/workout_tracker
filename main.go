package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joro550/workout_tracker/list"
	"github.com/joro550/workout_tracker/migrations"
	"github.com/joro550/workout_tracker/profile"
	"github.com/joro550/workout_tracker/task"
	"github.com/joro550/workout_tracker/users"
)

func main() {
	log.Println("🟩 Starting")

	db, err := sql.Open("sqlite3", "file:workout_tracker.db")
	if err != nil {
		log.Fatalln("💥 Could not connect to database", err)
	}

	for _, migration := range migrations.Migrations {
		_, err = db.Exec(migration.Run())
		if err != nil {
			log.Fatalln("💥 Migration failed", err)
		}
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	hello

	users.RegisterUserEndpoints(router, db)
	profile.RegisterProfileEndpoints(router, db)
	list.RegisterListEndpoints(router, db)
	task.RegisterTaskEndpoints(router, db)

	log.Println("👂 listening")

	http.ListenAndServe(":8080", router)
}
