package media

type VideoTrack interface {
	Identifier() string      // unique id
	Play() chan int64        // Start playing the track
	Pause() error            // pause playback (CurrentSegment stays the same)
	Stop() error             // pause playback and reset counter
	IsReady() bool           // Is the track ready to be played now?
	IsDone() bool            // Is the track done playing?
	CurrentSegment() Segment // Describe the next Segment to append
	Class() string
}

type Segment struct {
	URL           string
	Duration      float64
	Title         string
	Discontinuity bool
}

func (s Segment) Params() []interface{} {
	return []interface{}{s.URL, s.Duration, s.Title}
}

func BlankTrack(trackClass string) VideoTrack {
	switch trackClass {
	case "saved":
		return SavedTrack{}
	default:
		return SavedTrack{}
	}
}
