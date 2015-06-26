package datastore

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"strconv"
	"time"
)

// struct Channel describes a broadcast channel
type Channel struct {
	Identifier      string // Identifier is the unique part of the redis key
	PlaybackCounter int64
	mediaPlaylist   *m3u8.MediaPlaylist
}

// function GetChannel retrieves Channel information from the datastore
// given a channel identifier. It returns the channel and any error
func GetChannel(channelId string) (rv *Channel, err error) {
	rv = &Channel{}
	rv.Identifier = channelId
	rv.PlaybackCounter, err = client.IncrBy(rv.PlaybackCounterKey(), 0).Result()
	rv.mediaPlaylist, _ = m3u8.NewMediaPlaylist(1000, 1000)
	return
}

func (c *Channel) PlaylistData() string {
	return c.mediaPlaylist.Encode().String()
}

func (c Channel) CurrentItem() (PlaylistItem, error) {
	jsonString, err := client.LIndex(c.PlaybackQueueKey(), 0).Result()
	if err != nil {
		return PlaylistItem{}, err
	}

	item := PlaylistItem{}
	if err := json.Unmarshal([]byte(jsonString), &item); err != nil {
		return PlaylistItem{}, err
	}

	return item, nil
}

func (c Channel) PushItem(i *PlaylistItem) error {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return err
	}

	err = client.RPush(c.PlaybackQueueKey(), string(jsonBytes)).Err()
	return err
}

// function AdvanceCounter advances the playback counter, and updates the datastore
func (c *Channel) AdvanceCounter() error {
	currentItem, err := c.CurrentItem()
	if err != nil {
		return err
	}

	if currentItem.StartAt > c.PlaybackCounter {
		c.PlaybackCounter = currentItem.StartAt
	}
	videoFile := fmt.Sprintf(currentItem.URLFormat, c.PlaybackCounter)
	if err := c.mediaPlaylist.Append(videoFile, 5.0, ""); err != nil {
		log.Println("Error appending item to playlist:", err)
		c.mediaPlaylist.Slide(videoFile, 5.0, "")
	}
	c.PlaybackCounter = c.PlaybackCounter + 1

	if c.PlaybackCounter > currentItem.EndAt {
		c.ResetCounter()
		if !currentItem.Loop {
			client.LPop(c.PlaybackQueueKey())
		}
	}

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

// function PlaybackQueueKey returns the data store key for the playback queue
func (c Channel) PlaybackQueueKey() string {
	return "broadcast-channel-" + c.Identifier + "-queue"
}

// function PlaybackCounterKey returns the data store key for the playback counter
func (c Channel) PlaybackCounterKey() string {
	return "broadcast-channel-" + c.Identifier + "-counter"
}
