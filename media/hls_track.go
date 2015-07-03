package media

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/likexian/simplejson-go"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type HLSTrack struct {
	Id              string
	PlaylistURL     string
	PlaybackCounter int64
	tickChannel     chan int64
	playlist        *m3u8.MediaPlaylist
}

func (hls *HLSTrack) PlayFrom(npc int64) chan int64 {
	if hls.tickChannel == nil {
		hls.PlaybackCounter = 0
		hls.tickChannel = make(chan int64, 10)
		go hls.StreamPlaylist()
	}

	return hls.tickChannel
}

func (hls *HLSTrack) StreamPlaylist() {
	for {
		resp, err := http.Get(hls.PlaylistURL)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close()
		pl, _ := m3u8.NewMediaPlaylist(1000, 1000)
		pl.DecodeFrom(resp.Body, false)
		hls.playlist = pl
		numSegments := 0
		for cont := true; cont; {
			if pl.Segments[numSegments] != nil {
				numSegments = numSegments + 1
			} else {
				cont = false
			}
		}
		trackSequence := uint64(int(pl.SeqNo) + numSegments - 1)
		if trackSequence > uint64(hls.PlaybackCounter) {
			hls.PlaybackCounter = int64(trackSequence)
			hls.tickChannel <- int64(trackSequence)
		}

		time.Sleep(1 * time.Second)
	}
}

func (hls HLSTrack) Pause() error {
	return nil
}

func (hls HLSTrack) Stop() error {
	return nil
}

func (hls HLSTrack) Identifier() string {
	return hls.Id
}

func (hls HLSTrack) Class() string {
	return "hls"
}

func (hls HLSTrack) IsReady() bool {
	return true
}

func (hls HLSTrack) IsDone() bool {
	return false
}

func (hls HLSTrack) Serialize() (string, error) {
	bytes, err := json.Marshal(hls)
	return string(bytes), err
}

func (hls *HLSTrack) Load(json string) error {
	obj, err := simplejson.Loads(json)
	if err != nil {
		return err
	}

	hls.Id, _ = obj.Get("Id").String()
	hls.PlaylistURL, _ = obj.Get("PlaylistURL").String()
	hls.PlaybackCounter, _ = obj.Get("PlaybackCounter").Int64()

	return nil
}

func (hls *HLSTrack) SegmentNumber(seg int64) Segment {
	actualSegmentNumber := seg - int64(hls.playlist.SeqNo)
	myseg := hls.playlist.Segments[actualSegmentNumber]
	if myseg == nil {
		fmt.Println(hls.playlist.Segments)
		return Segment{}
	}
	providedURI := myseg.URI
	if !strings.Contains(providedURI, "http://") && !strings.Contains(providedURI, "https://") {
		basePath := filepath.Dir(hls.PlaylistURL)
		providedURI = filepath.Join(basePath, providedURI)
	}

	return Segment{
		Title:         myseg.Title,
		URL:           providedURI,
		Duration:      float64(myseg.Duration),
		Discontinuity: myseg.Discontinuity,
	}
}
