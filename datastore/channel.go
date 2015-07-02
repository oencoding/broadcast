package datastore

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/omarqazi/broadcast/media"
	"log"
	"strconv"
	"time"
)

// struct Channel describes a broadcast channel
type Channel struct {
	Identifier    string // Identifier is the unique part of the redis key
	mediaPlaylist *m3u8.MediaPlaylist
}

// function GetChannel retrieves Channel information from the datastore
// given a channel identifier. It returns the channel and any error
func GetChannel(channelId string) (rv *Channel, err error) {
	rv = &Channel{}
	rv.Identifier = channelId
	rv.mediaPlaylist, _ = m3u8.NewMediaPlaylist(1000, 1000)
	return
}

func (c *Channel) PlaylistData() string {
	return c.mediaPlaylist.Encode().String()
}

func (c Channel) CurrentItem() (media.VideoTrack, error) {
	trackId, err := client.LIndex(c.PlaybackQueueKey(), 0).Result()
	if err != nil {
		return media.BlankTrack(""), err
	}

	trackType, err := client.Get(trackId + "-class").Result()
	if err != nil {
		return media.BlankTrack(""), err
	}

	item := media.BlankTrack(trackType)
	trackData, err := client.Get("track-" + trackId + "-data").Result()
	if err != nil {
		return media.BlankTrack(""), err
	}

	if err := json.Unmarshal([]byte(trackData), &item); err != nil {
		return media.BlankTrack(""), err
	}

	return item, nil
}

func (c Channel) PushItem(i VideoTrack) (err error) {
	trackId := i.Identifier()
	client.Set(trackId+"-class", i.Class())
	if jsonBytes, err := json.Marshal(i); err != nil {
		return err
	}
	client.Set("track-"+trackId+"-data", string(jsonBytes))
	err = client.RPush(c.PlaybackQueueKey(), i.Identifier()).Err()
	return err
}

func (c *Channel) Play() error {
	for {
		currentItem, err := c.CurrentItem()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		playback := currentItem.Play()
		timeout := time.Tick(10 * time.Second)

		select {
		case <-playback:
			segment := currentItem.CurrentSegment()
			params := segment.Params()
			if err := c.mediaPlaylist.Append(params...); err != nil {
				c.mediaPlaylist.Slide(params...)
			}
			if segment.Discontinuity {
				if err := c.mediaPlaylist.SetDiscontinuity(); err != nil {
					log.Println("Error setting discontinuity:", err)
				}
			}
		case <-timeout:
			segment := currentItem.CurrentSegment()
			params := segment.Params()
			if err := c.mediaPlaylist.Append(params...); err != nil {
				c.mediaPlaylist.Slide(params...)
			}
		}
	}
	return nil
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
