package profile

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/layouts"
	"github.com/joro550/workout_tracker/list"
	profile_pages "github.com/joro550/workout_tracker/profile/pages"
	"github.com/joro550/workout_tracker/users"
)

func RegisterProfileEndpoints(mux *chi.Mux, db *sql.DB) {
	listRepo := list.NewListRepository(db)

	mux.Group(func(r chi.Router) {
		r.Use(users.AuthCtx)
		r.Get("/profile", profile(listRepo))
	})
}

func profile(repo list.ListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		log.Println("got user", user.Id)

		view := profile_pages.Profile()
		page := layouts.Layout(view)
		page.Render(r.Context(), w)
	}
}
