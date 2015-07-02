package media

import (
	"time"
)

// Struct SavedTrack plays back a pre-recorded track split
// into equal length sequence files. This track allows broadcasting
// of segments generated via mediafilesegmenter or similar
type SavedTrack struct {
	Identifier      string
	URLFormat       string
	StartAt         int64
	EndAt           int64
	Loop            bool
	LoopUntil       *time.Time
	TargetDuration  float64
	tickChannel     chan int64
	PlaybackCounter int64
}

func (st *SavedTrack) Play() chan int64 {
	if st.tickChannel == nil {
		st.tickChannel = make(chan int64, 10)
		st.PlaybackCounter = 0
		go st.AdvanceEvery(time.Duration(st.TargetDuration * time.Second))
	}

	return st.tickChannel
}

func Pause() error {
	return nil
}

func Stop() error {
	return nil
}

func (st *SavedTrack) AdvanceEvery(d time.Duration) {
	for ; true; time.Sleep(d) {
		st.tickChannel <- st.PlaybackCounter
		st.PlaybackCounter = st.PlaybackCounter + 1
	}
}

func (st SavedTrack) IsReady() bool {
	return true
}

func (st SavedTrack) IsDone() bool {
	return st.PlaybackCounter > st.EndAt
}

func (st SavedTrack) CurrentSegment() Segment {
	return Segment{
		URL:      fmt.Sprintf(st.URLFormat, st.PlaybackCounter),
		Duration: st.TargetDuration,
		Title:    "",
	}
}
