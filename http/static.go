package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StaticHandler ...
type StaticHandler struct {
	*mux.Router
}

// NewStaticHandler ...
func NewStaticHandler() *StaticHandler {
	h := &StaticHandler{
		Router: mux.NewRouter(),
	}

	h.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return h
}
