package main

import (
	"log"
	"os"
	"os/signal"
	"pxy/http"
	"pxy/http/websockets"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	publishURL      = "rtmp://global-live.mux.com:5222/app"
)

var (
	subprotocols = []string{"streamKey"}
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	livestreamPool := websockets.NewLivestreamPool(
		readBufferSize,
		writeBufferSize,
		subprotocols,
		publishURL,
	)

	streamHandler := http.NewStreamHandler(livestreamPool)
	staticHandler := http.NewStaticHandler()

	handler := http.Handler{
		StreamHandler: streamHandler,
		StaticHandler: staticHandler,
	}

	server := http.Server{Handler: &handler, Addr: ":" + port}
	err := server.Open()
	if err != nil {
		log.Fatalln("Failed to start server: %w", err)
	} else {
		log.Println("Server is running.")
	}

	// Block until an interrupt signal is received to keep server alive.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println("Got signal:", s)
}
