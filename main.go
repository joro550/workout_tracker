package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joro550/workout_tracker/endpoints"
	"github.com/joro550/workout_tracker/migrations"
	"github.com/joro550/workout_tracker/templates/layouts"
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

	endpoints.Init(router, db)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := layouts.Main()
		component.Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", router)
}
