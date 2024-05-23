package endpoints

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/repositories"
	"github.com/joro550/workout_tracker/templates/layouts"
	"github.com/joro550/workout_tracker/templates/pages"
)

func RegisterUserEndpoints(mux *chi.Mux, db *sql.DB) {
	userRepo := repositories.NewUserRepo(db)

	mux.Route("/user", func(router chi.Router) {
		log.Println("💃 Registering user endpoint")
		router.Get("/login", login)
		router.Post("/login", loginUser)

		router.Get("/register", register)
		router.Post("/register", registerUser(userRepo))

		router.Group(func(r chi.Router) {
			r.Use(AuthCtx)
			r.Get("/profile", slash)
		})
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	login := pages.Login(false)
	page := layouts.Layout(login)
	page.Render(r.Context(), w)
}

func register(w http.ResponseWriter, r *http.Request) {
	registerView := pages.Register(false)
	layout := layouts.Layout(registerView)
	layout.Render(r.Context(), w)
}

func registerUser(repo repositories.UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		exists, err := repo.UserExists(username)
		if err != nil {
			log.Println("🤔 Query failed", err)
			register(w, r)
			return
		}

		if exists {
			registerView := pages.Register(true)
			layout := layouts.Layout(registerView)
			layout.Render(r.Context(), w)
			return
		}

		w.Write([]byte(username))
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func slash(w http.ResponseWriter, r *http.Request) {
	log.Println("slash")
	w.Write([]byte("Hello"))
}

func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔐 auth")

		_, e := r.Cookie("auth")
		if e != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", "thing")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
