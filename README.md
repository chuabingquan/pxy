# pxy
`pxy` is a Go server that routes incoming livestream data from websockets to an external RTMP endpoint.

__This project is a work in progress, I'll update more once I'm done with my sleep.__

## Context
For a side-project of mine, I've to broadcast livestreams through an external service that uses the RTMP protocol. Flutter, and all the web browsers out there do not support said protocol. Therefore, I built `pxy` to proxy the livestreams from such clients (via websockets) to the broadcasting RTMP servers. Since RTMP is still widely used in the video streaming industry, I thought amateurs like me could benefit from an implementation like `pxy` for our side-projects and such.

## Status
Fundamentally, `pxy` works well so far. However, there are probably still bugs that needs be ironed out. If you do find any, feel free to open an issue or make a pull request. Meanwhile, `pxy` could be used as a reference for how a websocket-RTMP proxy could be built. Using it in production now is a really bad idea.

## Try it Out
1. Firstly, install FFmpeg on your machine (just Google it).

2. Proceed to clone the project with the command below
```bash
git clone https://github.com/chuabingquan/pxy.git && cd pxy-master/
```
3. Update your RTMP endpoint address under the constants in `cmd/pxy/main.go`.
```go
const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	publishURL      = "rtmp://global-live.mux.com:5222/app" // This one here.
)
```
4. Execute the following command to build and run the `pxy` server
```bash
go run cmd/pxy/main.go
```
5. Access `http://localhost:8080` in your browser and supply your stream name/stream key. `pxy` will append it behind the specified RTMP endpoint address that's mentioned Step 3 `(e.g. rtmp://global-live.mux.com:5222/app/{YOUR_STREAM_NAME/KEY})`