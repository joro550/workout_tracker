package list

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/layouts"
	list_pages "github.com/joro550/workout_tracker/list/pages"
	"github.com/joro550/workout_tracker/users"
)

func RegisterListEndpoints(mux *chi.Mux, db *sql.DB) {
	repo := NewListRepository(db)

	mux.Group(func(r chi.Router) {
		r.Use(users.AuthCtx)

		r.Route("/list", func(chi chi.Router) {
			chi.Get("/add", addView)
			chi.Get("/cards", cards(repo))
			chi.Get("/edit", editList(repo))
			chi.Post("/add", addViewPost(repo))

			chi.Delete("/{id}", deleteList(repo))
		})
	})
}

func addView(w http.ResponseWriter, r *http.Request) {
	view := list_pages.AddList()
	authedLayout := layouts.Authed(view)
	page := layouts.Layout(authedLayout)

	page.Render(r.Context(), w)
}

func editList(repo ListRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		id, _ := strconv.Atoi(r.PathValue("id"))

		list, _ := repo.GetList(id, user.Id)

		view := list_pages.EditList(list)
		authedLayout := layouts.Authed(view)
		page := layouts.Layout(authedLayout)

		page.Render(r.Context(), w)
	}
}

func cards(repo ListRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		lists, err := repo.GetAllLists(user.Id)
		if err != nil {
			log.Println("could not get lists", err)
		}
		var cards []list_pages.CardModel
		for _, list := range lists {
			cards = append(cards, list_pages.CardModel{
				Id:          list.Id,
				Name:        list.Name,
				Description: list.Description,
			})
		}

		cardView := list_pages.Cards(cards)
		cardView.Render(r.Context(), w)
	}
}

func addViewPost(repo ListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		model := list_pages.CreateCardModelFromRequest(r)
		user := r.Context().Value("user").(users.User)

		_, err := repo.CreateList(List{Name: model.Name, Description: model.Description, UserId: user.Id})
		if err != nil {
			log.Println("💥 Could not create list ", err)
			return
		} else {
			log.Println("Created list for user", user.Id)
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func editViewPost(repo ListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		model, err := list_pages.EditCardModelFromRequest(r)
		if err != nil {
			log.Print("💥 edit card model could not be parsed from request", err)
			return
		}
		user := r.Context().Value("user").(users.User)

		_, err = repo.UpdateList(List{Id: model.Id, Name: model.Name, Description: model.Description, UserId: user.Id})
		if err != nil {
			log.Println("💥 Could not create list ", err)
			return
		} else {
			log.Println("Created list for user", user.Id)
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func deleteList(repo ListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(users.User)
		id, _ := strconv.Atoi(r.PathValue("id"))
		_ = repo.DeleteList(id, user.Id)

		w.Header().Set("HX-Trigger", "list-updated")
		w.WriteHeader(http.StatusOK)
	}
}
