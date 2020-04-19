package websockets

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// LivestreamPool ...
type LivestreamPool struct {
	connections map[string]*livestreamConnection
	upgrader    websocket.Upgrader
	lock        sync.RWMutex
}

// NewLivestreamPool ...
func NewLivestreamPool(readBufferSize, writeBufferSize int, subprotocols []string) *LivestreamPool {
	return &LivestreamPool{
		connections: make(map[string]*livestreamConnection),
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

	newStreamer := newLivestreamConnection(sessionID, userID, streamKey, conn)

	if _, exists := lp.connections[sessionID]; !exists {
		lp.connections[sessionID] = newStreamer
	} else {
		lp.connections[sessionID].Websocket.Close()
		lp.connections[sessionID] = newStreamer
	}

	go func() {
		for {
			messageType, payload, err := newStreamer.Websocket.ReadMessage()
			if err != nil || messageType == websocket.CloseMessage {
				lp.lock.Lock()
				defer lp.lock.Unlock()

				newStreamer.Websocket.Close()
				if existingStreamer, exists := lp.connections[sessionID]; exists {
					if existingStreamer.ID == newStreamer.ID {
						delete(lp.connections, sessionID)
					}
				}
				break
			}

			log.Println(payload)
		}
	}()

	return nil
}

// CloseUpdates ...
func (lp *LivestreamPool) CloseUpdates(sessionID string) error {
	lp.lock.Lock()
	defer lp.lock.Unlock()

	if existingStreamer, exists := lp.connections[sessionID]; exists {
		existingStreamer.Websocket.Close()
		delete(lp.connections, sessionID)
	} else {
		return errors.New("Session to close livestream for does not exist")
	}

	return nil
}
