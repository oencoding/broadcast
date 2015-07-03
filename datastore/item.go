package datastore

import (
	"github.com/omarqazi/broadcast/media"
	"time"
)

// struct PlaylistItem describes a video playback instruction
// in a channel's playback queue
type PlaylistItem struct {
	TrackId   string // identifier of the track
	Loop      bool   // whether the video should loop indefinitely
	LoopUntil *time.Time
}

func (pi PlaylistItem) VideoTrack() media.VideoTrack {
	return GetVideoTrack(pi.TrackId)
}
