package http

import "net/http"

// LivestreamUpdates ...
type LivestreamUpdates interface {
	RegisterStreamer(sessionID, userID, streamKey string, w http.ResponseWriter, r *http.Request) error
	CloseUpdates(sessionID string) error
}
