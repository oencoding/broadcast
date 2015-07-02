package media

type VideoTrack interface {
	Play() chan int64 // Start playing the track
	Pause() error
	Stop() error
	IsReady() bool
	IsDone() bool
	CurrentSegment() Segment
}

type Segment struct {
	URL      string
	Duration float64
	Title    string
}
