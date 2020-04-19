package http

import (
	"net/http"
	"strings"
)

// Handler implements the http.Handler interface and serves as the main handler for the server
// by redirecting requests to sub-handlers.
type Handler struct {
	StreamHandler *StreamHandler
	StaticHandler *StaticHandler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")

	if len(urlSegments) > 2 && urlSegments[0] == "api" {
		resourceName := urlSegments[3]

		switch resourceName {
		case "stream":
			h.StreamHandler.ServeHTTP(w, r)
			break
		default:
			http.NotFound(w, r)
			break
		}

		return
	}

	h.StaticHandler.ServeHTTP(w, r)
}
