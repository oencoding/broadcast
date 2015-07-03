package media

import (
	"encoding/json"
	"fmt"
	"github.com/likexian/simplejson-go"
	"time"
)

// Struct SavedTrack plays back a pre-recorded track split
// into equal length sequence files. This track allows broadcasting
// of segments generated via mediafilesegmenter or similar
type SavedTrack struct {
	Id              string
	URLFormat       string
	StartAt         int64
	EndAt           int64
	TargetDuration  float64
	tickChannel     chan int64
	PlaybackCounter int64
}

func (st *SavedTrack) PlayFrom(npc int64) chan int64 {
	if st.tickChannel == nil {
		st.PlaybackCounter = npc
		st.tickChannel = make(chan int64, 10)
		go st.AdvanceEvery(time.Duration(st.TargetDuration * float64(time.Second)))
	}

	return st.tickChannel
}

func (st SavedTrack) Pause() error {
	return nil
}

func (st SavedTrack) Stop() error {
	return nil
}

func (st *SavedTrack) AdvanceEvery(d time.Duration) {
	for !st.IsDone() {
		st.tickChannel <- st.PlaybackCounter
		st.PlaybackCounter = st.PlaybackCounter + 1
		time.Sleep(d)
	}
}

func (st SavedTrack) Identifier() string {
	return st.Id
}

func (st SavedTrack) Class() string {
	return "saved"
}

func (st SavedTrack) IsReady() bool {
	return true
}

func (st *SavedTrack) IsDone() bool {
	return st.PlaybackCounter > st.EndAt
}

func (st SavedTrack) Serialize() (string, error) {
	bytes, err := json.Marshal(st)
	return string(bytes), err
}

func (rv *SavedTrack) Load(json string) error {
	obj, err := simplejson.Loads(json)
	if err != nil {
		return err
	}

	rv.Id, _ = obj.Get("Id").String()
	rv.URLFormat, _ = obj.Get("URLFormat").String()
	rv.StartAt, _ = obj.Get("StartAt").Int64()
	rv.EndAt, _ = obj.Get("EndAt").Int64()
	rv.TargetDuration, _ = obj.Get("TargetDuration").Float64()
	rv.PlaybackCounter, _ = obj.Get("PlaybackCounter").Int64()

	return nil
}

func (st SavedTrack) SegmentNumber(seg int64) Segment {
	return Segment{
		URL:           fmt.Sprintf(st.URLFormat, seg),
		Duration:      st.TargetDuration,
		Title:         "",
		Discontinuity: (seg == st.StartAt),
	}
}
