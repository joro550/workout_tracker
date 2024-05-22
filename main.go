package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joro550/workout_tracker/endpoints"
	"github.com/joro550/workout_tracker/templates/layouts"
)

func main() {
	router := chi.NewRouter()
	endpoints.Init(router)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := layouts.Main()
		component.Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", router)
}
