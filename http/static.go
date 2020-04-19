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

	fs := http.FileServer(http.Dir("static"))
	h.Handle("/", http.StripPrefix("/", fs))

	return h
}
