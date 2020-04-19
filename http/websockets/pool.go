package websockets

import (
	"errors"
	"fmt"
	"net/http"
	"pxy"
	"pxy/rtmp"
	"sync"

	"github.com/gorilla/websocket"
)

type livestream struct {
	Streamer   *livestreamConnection
	RTMPClient pxy.BroadcastService
}

// LivestreamPool ...
type LivestreamPool struct {
	publishURL  string
	connections map[string]*livestream
	upgrader    websocket.Upgrader
	lock        sync.RWMutex
}

// NewLivestreamPool ...
func NewLivestreamPool(readBufferSize, writeBufferSize int, subprotocols []string, publishURL string) *LivestreamPool {
	return &LivestreamPool{
		publishURL:  publishURL,
		connections: make(map[string]*livestream),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
			Subprotocols:    subprotocols,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		lock: sync.RWMutex{},
	}
}

// RegisterStreamer ...
func (lp *LivestreamPool) RegisterStreamer(sessionID, userID, streamKey string, w http.ResponseWriter, r *http.Request) error {
	lp.lock.Lock()
	defer lp.lock.Unlock()

	conn, err := lp.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("Failed to upgrade HTTP connection to a websocket: %w", err)
	}

	rtmp := rtmp.NewRTMPClient(lp.publishURL, streamKey)
	err = rtmp.StartBroadcast()
	if err != nil {
		return fmt.Errorf("Failed to start RTMP broadcast: %w", err)
	}

	newStreamer := newLivestreamConnection(sessionID, userID, streamKey, conn)

	livestream := &livestream{
		Streamer:   newStreamer,
		RTMPClient: rtmp,
	}

	if _, exists := lp.connections[sessionID]; !exists {
		lp.connections[sessionID] = livestream
	} else {
		lp.connections[sessionID].RTMPClient.StopBroadcast()
		lp.connections[sessionID].Streamer.Websocket.Close()
		lp.connections[sessionID] = livestream
	}

	go func() {
		for {
			messageType, payload, err := newStreamer.Websocket.ReadMessage()
			if err != nil || messageType == websocket.CloseMessage {
				lp.lock.Lock()
				defer lp.lock.Unlock()

				rtmp.StopBroadcast()
				newStreamer.Websocket.Close()
				if existingLivestream, exists := lp.connections[sessionID]; exists {
					if existingLivestream.Streamer.ID == newStreamer.ID {
						delete(lp.connections, sessionID)
					}
				}
				break
			}

			rtmp.PipeToBroadcast(&payload)
		}
	}()

	return nil
}

// RemoveStreamer ...
func (lp *LivestreamPool) RemoveStreamer(sessionID string) error {
	lp.lock.Lock()
	defer lp.lock.Unlock()

	if existingLivestream, exists := lp.connections[sessionID]; exists {
		existingLivestream.RTMPClient.StopBroadcast()
		existingLivestream.Streamer.Websocket.Close()
		delete(lp.connections, sessionID)
	} else {
		return errors.New("Session to close livestream for does not exist")
	}

	return nil
}

// StreamerIsRegisteredForSession ...
func (lp *LivestreamPool) StreamerIsRegisteredForSession(sessionID string) bool {
	lp.lock.Lock()
	defer lp.lock.Unlock()

	_, exists := lp.connections[sessionID]

	return exists
}
