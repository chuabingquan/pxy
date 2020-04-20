# pxy
`pxy` is a Go server that routes incoming livestream data from websockets to an external RTMP endpoint.

__This project is a work in progress, I'll update it more with time.__

## Context
For a side-project of mine, I've to broadcast live streams through an external service that uses the RTMP protocol. Unfortunately, my front-ends (Flutter and all web browsers out there) do not support the RTMP protocol. Therefore, I built `pxy` to proxy the live streams from such clients (via websockets) to the broadcasting RTMP servers. Since RTMP is still widely used in the video streaming industry, I thought amateurs like myself could benefit from an implementation like `pxy` for our own live streaming side-projects.

## Status
Fundamentally, `pxy` works well so far. However, there are probably still bugs that needs be ironed out. If you do find any, feel free to open an issue or make a pull request. Meanwhile, you could use `pxy` as a reference for implementing your own websocket-RTMP proxy. Alternatively, you could clone this project and modify it to suit your needs.

## Try it Out
`pxy` can be setup in two ways, with, or without Docker. Head to the respective sections after completing the preliminaries for more instructions on your preferred way to setup.

### Preliminaries
Clone and navigate to the repository's root.
```bash
git clone https://github.com/chuabingquan/pxy.git && cd pxy-master/
```

Update your RTMP endpoint address under the constants in `cmd/pxy/main.go`.
```go
const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	publishURL      = "rtmp://global-live.mux.com:5222/app" // This one here.
)
```

### Setup pxy with Docker (Recommended)
_Before proceeding, please ensure that Docker is already installed on your machine._

Start by building the `pxy` docker image.
```bash
docker build -t pxy .
```

Proceed to run a container that's based on the built `pxy` docker image.
```bash
docker run --rm -p 8080:8080 pxy
```


### Setup pxy without Docker
_Before proceeding, please install FFmpeg on your machine (please Google how to)._

Build and run the `pxy` server.
```bash
go run cmd/pxy/main.go
```

### Final Steps
Once your `pyx` instance is up and running, access `http://localhost:8080` from your browser (preferably Chrome) and supply your stream name/stream key. 

On the backend, `pxy` will append the given stream name/stream key behind the RTMP endpoint address that's defined in the preliminaries section of this README `(e.g. rtmp://global-live.mux.com:5222/app/<YOUR_STREAM_NAME>)`.

At this point, your RTMP endpoint should receive the proxied live stream from `pxy`. Playback (viewing) of the live stream is dependent on the tools/services that's hosting your RTMP endpoint.

## References
- [The state of going live from a browser](https://mux.com/blog/the-state-of-going-live-from-a-browser/)
- [Streaming to Facebook Live from a canvas](https://github.com/fbsamples/Canvas-Streaming-Example/blob/master/README.md)