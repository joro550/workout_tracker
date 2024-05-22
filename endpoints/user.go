package endpoints

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joro550/workout_tracker/templates/layouts"
	"github.com/joro550/workout_tracker/templates/pages"
)

func Init(mux *chi.Mux) {
	mux.Route("/user", func(router chi.Router) {
		router.Get("/login", login)
		router.Post("/login", loginUser)
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	login := pages.Login()
	page := layouts.Layout(login)
	page.Render(r.Context(), w)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
