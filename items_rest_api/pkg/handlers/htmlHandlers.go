package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"github.com/danrusei/items-rest-api/pkg/model"
)

func (h *Handlers) htmlHandleList() http.HandlerFunc {
	var (
		init sync.Once
		tmpl *template.Template
		err  error
	)

	type ListItemPage struct {
		PageTitle string
		Items     []model.Item
	}

	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tmpl, err = template.ParseFiles(h.htmlFiles...)
		})
		retrieve, err := h.db.ListGoods()
		if err != nil {
			http.Error(w, "couldn't retrieve data", http.StatusInternalServerError)
		}

		data := ListItemPage{
			PageTitle: "Item Database",
			Items:     retrieve,
		}

		tmpl.Execute(w, data)
	}
}

func (h *Handlers) htmlHandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respItem struct {
			Name      string          `json:"name"`
			ExpDate   model.Timestamp `jsnon:"expdate"`
			ExpOpen   string          `json:"expopen"`
			Comment   string          `json:"comment"`
			TargetAge string          `json:"targetage"`
			IsOpen    string          `json:"isopen"`
			Opened    model.Timestamp `json:"opened"`
		}

		err := h.decode(w, r, &respItem)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}

		expopen, err := strconv.Atoi(respItem.ExpOpen)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}

		isopen, err := strconv.ParseBool(respItem.IsOpen)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}

		item := model.Item{
			Good: model.Good{
				Name:    respItem.Name,
				ExpDate: respItem.ExpDate,
				ExpOpen: expopen,
				Comment: respItem.Comment,
			},
			TargetAge: respItem.TargetAge,
			IsOpen:    isopen,
			Opened:    respItem.Opened,
		}

		data, err := h.db.AddGood(item)
		if err != nil {
			h.respond(w, r, data, http.StatusBadRequest)
		}

		h.respond(w, r, data, http.StatusOK)
	}
}
