package datastore

import (
	"strconv"
	"time"
)

// struct Channel describes a broadcast channel
type Channel struct {
	Identifier      string // Identifier is the unique part of the redis key
	PlaybackCounter int64
}

// function GetChannel retrieves Channel information from the datastore
// given a channel identifier. It returns the channel and any error
func GetChannel(channelId string) (rv *Channel, err error) {
	rv = &Channel{}
	rv.Identifier = channelId
	rv.PlaybackCounter, err = client.IncrBy(rv.PlaybackCounterKey(), 0).Result()
	return
}

// function AdvanceCounter advances the playback counter, and updates the datastore
func (c *Channel) AdvanceCounter() error {
	c.PlaybackCounter = c.PlaybackCounter + 1
	return c.SaveCounter()
}

// function ResetCounter resets the playback counter to 0, and updates the datastore
func (c *Channel) ResetCounter() error {
	c.PlaybackCounter = 0
	return c.SaveCounter()
}

// function AdvanceEvery increments the counter every duration d
// the second argument should be a channel that will cancel the operation on receive
func (c *Channel) AdvanceEvery(d time.Duration, cancel chan int) {
	ticker := time.Tick(d)
	c.AdvanceCounter()

	for {
		select {
		case <-cancel:
			return
		case <-ticker:
			c.AdvanceCounter()
		}
	}

	return
}

// function SaveCounter saves the current PlaybackCounter in the data store
func (c Channel) SaveCounter() error {
	return client.Set(c.PlaybackCounterKey(), strconv.FormatInt(c.PlaybackCounter, 10), 0).Err()
}

// function PlaybackCounterKey returns the data store key for the playback counter
func (c Channel) PlaybackCounterKey() string {
	return "broadcast-channel-" + c.Identifier + "-counter"
}
