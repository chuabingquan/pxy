package rtmp

import (
	"fmt"
	"os"
	"os/exec"
)

// Client ...
type Client struct {
	RtmpURL    string
	StreamKey  string
	streamPipe *chan *[]byte
	ffmpeg     *exec.Cmd
}

var ffmpegArgs = []string{
	"-i", "-",
	"-vcodec", "copy",
	"-f", "flv",
}

// NewRTMPClient ...
func NewRTMPClient(rtmpURL, streamKey string) *Client {
	streamPipe := make(chan *[]byte)
	return &Client{
		RtmpURL:    rtmpURL,
		StreamKey:  streamKey,
		streamPipe: &streamPipe,
		ffmpeg:     exec.Command("ffmpeg", append(ffmpegArgs, rtmpURL+"/"+streamKey)...),
	}
}

// StartBroadcast ...
func (c *Client) StartBroadcast() error {
	c.ffmpeg.Stderr = os.Stderr

	ffmpegInput, err := c.ffmpeg.StdinPipe()
	if err != nil {
		return fmt.Errorf("Failed to get input pipe for FFmpeg process: %w", err)
	}

	err = c.ffmpeg.Start()
	if err != nil {
		ffmpegInput.Close()
		return fmt.Errorf("Failed to start FFmpeg process: %w", err)
	}

	go func() {
		defer ffmpegInput.Close()
		defer c.ffmpeg.Wait()

		for {
			stream := <-*c.streamPipe

			_, err := ffmpegInput.Write(*stream)
			if err != nil {
				break
			}
		}
	}()

	return nil
}

// StopBroadcast ...
func (c *Client) StopBroadcast() error {
	err := c.ffmpeg.Process.Kill()
	if err != nil {
		return fmt.Errorf("Failed to kill FFmpeg process: %w", err)
	}

	return nil
}

// PipeToBroadcast ...
func (c *Client) PipeToBroadcast(stream *[]byte) {
	*c.streamPipe <- stream
}
