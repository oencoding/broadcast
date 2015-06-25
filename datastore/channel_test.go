package datastore

import (
	"testing"
)

func TestGetChannel(t *testing.T) {
	if channel, err := GetChannel("test-channel"); err != nil {
		t.Fatal("Error getting channel:", channel)
	} else if channel.PlaybackCounter != 0 {
		t.Fatal("Error: expected playback counter of 0 but got", channel.PlaybackCounter)
	}
}
