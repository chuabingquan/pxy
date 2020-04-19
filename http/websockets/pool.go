package websockets

import "net/http"

// LivestreamUpdates ...
type LivestreamUpdates struct{}

// RegisterStreamer ...
func (lu *LivestreamUpdates) RegisterStreamer(sessionID, userID, streamKey string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// CloseUpdates ...
func (lu *LivestreamUpdates) CloseUpdates(sessionID string) error {
	return nil
}
