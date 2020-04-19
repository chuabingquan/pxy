package pxy

// BroadcastService ...
type BroadcastService interface {
	StartBroadcast() error
	StopBroadcast() error
	PipeToBroadcast(stream *[]byte)
}
