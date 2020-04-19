package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// StreamHandler ..
type StreamHandler struct {
	*mux.Router
	livestreamUpdates LivestreamUpdates
}

// NewStreamHandler ...
func NewStreamHandler(livestreamUpdates LivestreamUpdates) *StreamHandler {
	h := &StreamHandler{
		Router:            mux.NewRouter(),
		livestreamUpdates: livestreamUpdates,
	}

	h.HandleFunc("/api/v0/stream", h.handleCreateLivestream).Methods("GET")

	return h
}

func (sh *StreamHandler) handleCreateLivestream(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["sessionID"]
	userID := "some_user_id"

	streamKey, err := getWebsocketStreamKey(r)
	if err != nil {
		writeJSONMessage(w, http.StatusBadRequest, standardResponse{"Please supply a valid stream key"})
		return
	}

	err = sh.livestreamUpdates.RegisterStreamer(sessionID, userID, streamKey, w, r)
	if err != nil {
		writeJSONMessage(w, http.StatusInternalServerError, standardResponse{"Something unexpected went wrong"})
		return
	}
}

func getWebsocketStreamKey(r *http.Request) (string, error) {
	subprotocols := strings.Split(r.Header.Get("Sec-Websocket-Protocol"), ", ")
	if len(subprotocols) != 2 || subprotocols[0] != "streamKey" {
		return "", errors.New("Stream key not found in header")
	}
	return subprotocols[1], nil
}
