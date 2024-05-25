package list_pages

import "net/http"

type CardModel struct {
	Name        string
	Description string
	Id          int
}

type CreateCardModel struct {
	Name        string
	Description string
}

func CreateCardModelFromRequest(r *http.Request) *CreateCardModel {
	return &CreateCardModel{
		Name:        r.FormValue("Name"),
		Description: r.FormValue("Description"),
	}
}
