package list_pages

import (
	"net/http"
	"strconv"
)

type CardModel struct {
	Name        string
	Description string
	Id          int
}

type CreateCardModel struct {
	Name        string
	Description string
}

type EditCardModel struct {
	Name        string
	Description string
	Id          int
}

func CreateCardModelFromRequest(r *http.Request) *CreateCardModel {
	return &CreateCardModel{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}
}

func EditCardModelFromRequest(r *http.Request) (*EditCardModel, error) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, err
	}
	return &EditCardModel{
		Id:          id,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}, nil
}
