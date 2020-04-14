package handlers

import "net/http"

//CreateRoutes wire up the server endpoints
func (h *Handlers) CreateRoutes(mux *http.ServeMux) {

	//json handlers
	mux.HandleFunc("/json/", h.logging(h.jsonHandleList()))
	mux.HandleFunc("/json/add", h.logging(h.jsonHandleAdd()))
	mux.HandleFunc("/json/open", h.logging(h.jsonHandleOpen()))
	mux.HandleFunc("/json/del", h.logging(h.jsonHandleDelete()))

	//html handlers
	mux.HandleFunc("/", h.logging(h.htmlHandleList()))
	mux.HandleFunc("/add", h.logging(h.htmlHandleAdd()))
}
