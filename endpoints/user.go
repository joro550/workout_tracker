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
		router.Post("/login", loginUser(&userRepo))

		router.Get("/register", register(pages.RegisterModel{}))
		router.Post("/register", registerUser(&userRepo))
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	login := pages.Login(pages.LoginModel{})
	page := layouts.Layout(login)
	page.Render(r.Context(), w)
}

func register(model pages.RegisterModel) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		registerView := pages.Register(model)
		layout := layouts.Layout(registerView)
		layout.Render(r.Context(), w)
	}
}

func registerUser(repo *repositories.UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		exists, err := repo.UserExists(username)
		if err != nil {
			log.Println("🤔 Query failed", err)
			register(pages.RegisterModel{})(w, r)
			return
		}

		if exists {
			view := register(pages.RegisterModel{UserExists: true})
			view(w, r)
			return
		}

		encPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := repositories.User{
			Username: username,
			Password: string(encPass),
		}

		id, _ := repo.CreateUser(user)
		log.Println("👶 Id for new user was", id)
		token, _ := createUserCookie(id, &user)
		http.SetCookie(w, &http.Cookie{Value: token, Name: "jwt", Path: "/"})
		repo.CreateUser(user)
	}
}

func loginUser(repo *repositories.UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		exists, _ := repo.UserExists(username)
		if !exists {
			login(w, r)
			return
		}

		user, _ := repo.GetUser(username)
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		passwordCorrect := bcrypt.CompareHashAndPassword([]byte(user.Password), hash)
		if passwordCorrect != nil {
			login(w, r)
			return
		}

		token, _ := createUserCookie(user.Id, &user)
		http.SetCookie(w, &http.Cookie{Value: token, Name: "jwt", Path: "/"})
		http.Redirect(w, r, "profile", http.StatusFound)
	}
}

func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔐 auth")
		cookieString, err := r.Cookie("jwt")
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

		// this just wants to be a float64
		id := int(claims["id"].(float64))
		username := claims["username"].(string)

		ctx := context.WithValue(r.Context(), "user", repositories.User{Id: id, Username: username})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createUserCookie(id int, user *repositories.User) (string, error) {
	auth := createAuthToken()

	claims := map[string]interface{}{"id": id, "username": user.Username}
	_, tokenString, err := auth.Encode(claims)
	return tokenString, err
}

func createAuthToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("password"), nil)
}
