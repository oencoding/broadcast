package datastore

import "time"

// struct PlaylistItem describes a video playback instruction
// in a channel's playback queue
type PlaylistItem struct {
	TrackId   string // identifier of the track
	Loop      bool   // whether the video should loop indefinitely
	LoopUntil *time.Time
}
