package endpoints

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/repositories"
	"github.com/joro550/workout_tracker/templates/layouts"
	"github.com/joro550/workout_tracker/templates/pages"
)

func RegisterProfileEndpoints(mux *chi.Mux, db *sql.DB) {
	listRepo := repositories.NewListRepository(db)

	mux.Group(func(r chi.Router) {
		r.Use(AuthCtx)
		r.Get("/profile", profile(listRepo))
	})
}

func profile(repo repositories.ListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(repositories.User)
		log.Println("got user", user.Id)

		view := pages.Profile()
		page := layouts.Layout(view)
		page.Render(r.Context(), w)
	}
}
