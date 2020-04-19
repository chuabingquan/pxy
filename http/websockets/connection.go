package websockets

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type livestreamConnection struct {
	ID        string
	SessionID string
	UserID    string
	StreamKey string
	Websocket *websocket.Conn
}

func newLivestreamConnection(sessionID, userID, streamKey string, conn *websocket.Conn) *livestreamConnection {
	return &livestreamConnection{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		UserID:    userID,
		StreamKey: streamKey,
		Websocket: conn,
	}
}
