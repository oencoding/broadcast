package media

type VideoTrack interface {
	Identifier() string              // unique id
	PlayFrom(npc int64) chan int64   // Start playing the track
	Pause() error                    // pause playback (CurrentSegment stays the same)
	Stop() error                     // pause playback and reset counter
	IsReady() bool                   // Is the track ready to be played now?
	IsDone() bool                    // Is the track done playing?
	SegmentNumber(seg int64) Segment // Describe the next Segment to append
	Class() string
	Serialize() (string, error)
	Load(serialized string) error
}

type Segment struct {
	URL           string
	Duration      float64
	Title         string
	Discontinuity bool
}

func BlankTrack(trackClass string) VideoTrack {
	switch trackClass {
	case "saved":
		return &SavedTrack{}
	case "hls":
		return &HLSTrack{}
	default:
		return &SavedTrack{}
	}
}
