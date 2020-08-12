package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danrusei/items-rest-api/pkg/model"
)

func (h *Handlers) jsonHandleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.db.ListGoods()
		if err != nil {
			h.respond(w, r, err, http.StatusNotFound)
		}
		h.respond(w, r, data, http.StatusOK)
	}
}

func (h *Handlers) jsonHandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item model.Item
		err := h.decode(w, r, &item)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}

		data, err := h.db.AddGood(item)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}

		h.respond(w, r, data, http.StatusOK)
	}
}

func (h *Handlers) jsonHandleOpen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		b, err := strconv.ParseBool(r.FormValue("open"))
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}
		data, err := h.db.OpenState(id, b)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}
		h.respond(w, r, data, http.StatusOK)
	}
}

func (h *Handlers) jsonHandleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		data, err := h.db.DelGood(id)
		if err != nil {
			h.respond(w, r, err, http.StatusBadRequest)
		}
		h.respond(w, r, data, http.StatusOK)
	}
}

func (h *Handlers) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Could not encode in json", status)
		}
	}
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)

}
