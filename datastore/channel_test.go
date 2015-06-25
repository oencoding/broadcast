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

const incChannel = "increment-channel"

func TestMoveCounter(t *testing.T) {
	channel, err := GetChannel(incChannel)
	if err != nil {
		t.Fatal("Error getting channel:", err)
	}

	initialValue := channel.PlaybackCounter
	expectedValue := initialValue + 1
	if err := channel.AdvanceCounter(); err != nil {
		t.Fatal("Error advancing channel:", err)
	} else if nchan, _ := GetChannel(incChannel); nchan.PlaybackCounter != expectedValue {
		t.Error("Error: After incrementing channel expected value of", expectedValue, "but got", nchan.PlaybackCounter)
	}

	if err := channel.ResetCounter(); err != nil {
		t.Fatal("Error reseting counter:", err)
	} else if rchan, _ := GetChannel(incChannel); rchan.PlaybackCounter != 0 {
		t.Error("Error: After reseting counter expected value of 0 but got", rchan.PlaybackCounter)
	}
}
