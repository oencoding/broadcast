package datastore

// struct PlaylistItem describes a video playback instruction
// in a channel's playback queue
type PlaylistItem struct {
	URLFormat string // The format for the url, with sequence number as argument
	StartAt   int64  // The sequence number to start playing at
	EndAt     int64  // The sequence number to stop playing at
	Loop      bool   // whether the video should loop indefinitely
}
