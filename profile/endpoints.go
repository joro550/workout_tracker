package profile

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joro550/workout_tracker/layouts"
	"github.com/joro550/workout_tracker/list"
	list_pages "github.com/joro550/workout_tracker/list/pages"
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

		cardView := list_pages.UpdateableCards(cards)

		view := profile_pages.Profile(cardView)
		authedLayout := layouts.Authed(view)
		page := layouts.Layout(authedLayout)

		page.Render(r.Context(), w)
	}
}
