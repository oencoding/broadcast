package datastore

import (
	"testing"
	"time"
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

func TestAutoAdvance(t *testing.T) {
	channel, err := GetChannel("auto-channel")
	if err != nil {
		t.Fatal("Error getting channel:", err)
	}

	originalValue := channel.PlaybackCounter
	cancel := make(chan int, 1)

	go channel.AdvanceEvery(1*time.Millisecond, cancel)

	time.Sleep(10 * time.Millisecond)
	if channel.PlaybackCounter <= originalValue {
		t.Fatal("Error: playback counter did not automatically advance. Was", originalValue, "now", channel.PlaybackCounter)
	} else {
		originalValue = channel.PlaybackCounter
	}

	time.Sleep(10 * time.Millisecond)
	if channel.PlaybackCounter <= originalValue {
		t.Fatal("Error: playback counter did not automatically advance. Was", originalValue, "now", channel.PlaybackCounter)
	}

	cancel <- 1
	time.Sleep(10 * time.Millisecond) // give it a bit to cancel
	originalValue = channel.PlaybackCounter

	time.Sleep(10 * time.Millisecond)
	if channel.PlaybackCounter > originalValue {
		t.Fatal("Error: playback counter advanced after cancel. Was", originalValue, "then", channel.PlaybackCounter)
	}

}
