package endpoints

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joro550/workout_tracker/repositories"
	"github.com/joro550/workout_tracker/templates/layouts"
	"github.com/joro550/workout_tracker/templates/pages"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserEndpoints(mux *chi.Mux, db *sql.DB) {
	userRepo := repositories.NewUserRepo(db)

	mux.Route("/user", func(router chi.Router) {
		router.Get("/login", login)
		router.Post("/login", loginUser)

		router.Get("/register", register)
		router.Post("/register", registerUser(userRepo))
	})

	mux.Route("/profile", func(r chi.Router) {
		r.Use(AuthCtx)
		r.Get("/", profile)
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	login := pages.Login(false)
	page := layouts.Layout(login)
	page.Render(r.Context(), w)

	cookie, err := r.Cookie("auth")
	if err == nil {
		log.Println(cookie.Value)
	}
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

		encPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := repositories.User{
			Username: username,
			Password: string(encPass),
		}

		token, _ := createUserCookie(&user)
		http.SetCookie(w, &http.Cookie{Value: token, Name: "jwt", Path: "/"})
		repo.CreateUser(user)
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("setting cookie")
	username := r.FormValue("username")
	user := repositories.User{
		Username: username,
	}

	token, _ := createUserCookie(&user)
	http.SetCookie(w, &http.Cookie{Value: token, Name: "jwt", Path: "/"})
	w.Write([]byte("Hello"))
}

func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔐 auth")
		cookieString, err := r.Cookie("jwt")

		for _, cookie := range r.Cookies() {
			log.Println("name: [", cookie.Name, "]")
		}

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Println("no cookie", cookieString, err)
			return
		}

		auth, err := jwtauth.VerifyToken(createAuthToken(), cookieString.Value)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Println("something went wrong", cookieString, err)
			return
		}

		claims := auth.PrivateClaims()
		id := claims["id"]
		log.Println("userid: ", id)

		ctx := context.WithValue(r.Context(), "user", "thing")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createUserCookie(user *repositories.User) (string, error) {
	auth := createAuthToken()

	claims := map[string]interface{}{"id": user.Id, "username": user.Username}
	_, tokenString, err := auth.Encode(claims)
	return tokenString, err
}

func createAuthToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("password"), nil)
}
