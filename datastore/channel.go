package datastore

import (
	"encoding/json"
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
	rv.mediaPlaylist, _ = m3u8.NewMediaPlaylist(100, 100)
	return
}

// function PlaylistData() returns the current channel data in m3u8
// HTTP live streaming playlist format
func (c *Channel) PlaylistData() string {
	return c.mediaPlaylist.Encode().String()
}

// function CurrentItem() returns the currently playing VideoTrack
func (c Channel) CurrentItem() (media.VideoTrack, error) {
	trackId, err := client.LIndex(c.PlaybackQueueKey(), 0).Result()
	if err != nil {
		return media.BlankTrack(""), err
	}

	trackType, err := client.Get(trackId + "-class").Result()
	if err != nil {
		trackType = "saved"
	}

	item := media.BlankTrack(trackType)
	trackData, err := client.Get("track-" + trackId + "-data").Result()
	if err != nil {
		return media.BlankTrack(""), err
	}

	if err := item.Load(trackData); err != nil {
		return media.BlankTrack(""), err
	}

	return item, nil
}

// Function PushItem serializes a VideoTrack to storage, and queues
// it for playback at the end of the track queue
func (c Channel) PushItem(i media.VideoTrack) (err error) {
	trackId := i.Identifier()
	client.Set(trackId+"-class", i.Class(), 0)
	if jsonBytes, err := json.Marshal(i); err == nil {
		client.Set("track-"+trackId+"-data", string(jsonBytes), 0)
		err = client.RPush(c.PlaybackQueueKey(), i.Identifier()).Err()
	}

	return err
}

// Function Play starts the broadcast timer
func (c *Channel) Play() {
	for {
		if currentItem, err := c.CurrentItem(); err == nil {
			c.PlayTrack(currentItem)
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

// Function PlayTrack broadcasta a track on the channel until the
// track is finished
func (c *Channel) PlayTrack(currentItem media.VideoTrack) error {
	npc := int64(0)
	playback := currentItem.PlayFrom(c.GetPlaybackCounter())

	for {
		timeout := time.Tick(10 * time.Second)
		select {
		case npc = <-playback:
			c.BroadcastSegment(currentItem, npc, true)
		case <-timeout:
			c.BroadcastSegment(currentItem, npc, false)
		}

		c.SetPlaybackCounter(npc)

		if currentItem.IsDone() {
			return nil
		}
	}

	return nil
}

// function BroadcastSegment broadcasts a single segment on the channel
func (c *Channel) BroadcastSegment(v media.VideoTrack, pc int64, breaks bool) {
	segment := v.SegmentNumber(pc)
	err := c.mediaPlaylist.Append(segment.URL, segment.Duration, segment.Title)
	if err != nil {
		c.mediaPlaylist.Slide(segment.URL, segment.Duration, segment.Title)
	}

	if breaks && segment.Discontinuity {
		if err := c.mediaPlaylist.SetDiscontinuity(); err != nil {
			log.Println("Error setting discontinuity:", err)
		}
	}
}

func (c Channel) GetPlaybackCounter() int64 {
	cnt, err := client.Get(c.PlaybackCounterKey()).Result()
	if err != nil {
		return 0
	} else {
		rv, err := strconv.ParseInt(cnt, 10, 64)
		if err != nil {
			return 0
		} else {
			return rv
		}
	}
}

func (c Channel) SetPlaybackCounter(npc int64) error {
	err := client.Set(c.PlaybackCounterKey(), strconv.FormatInt(npc, 10), 0).Err()
	return err
}

// function PlaybackQueueKey returns the data store key for the playback queue
func (c Channel) PlaybackQueueKey() string {
	return "broadcast-channel-" + c.Identifier + "-queue"
}

// function PlaybackCounterKey returns the data store key for the playback counter
func (c Channel) PlaybackCounterKey() string {
	return "broadcast-channel-" + c.Identifier + "-counter"
}
