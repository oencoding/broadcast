package datastore

import "strconv"

// struct Channel describes a broadcast channel
type Channel struct {
	Identifier      string // Identifier is the unique part of the redis key
	PlaybackCounter int64
}

func GetChannel(channelId string) (rv Channel, err error) {
	rv.Identifier = channelId
	rv.PlaybackCounter, err = client.IncrBy(rv.PlaybackCounterKey(), 0).Result()
	return
}

func (c *Channel) AdvanceCounter() error {
	c.PlaybackCounter = c.PlaybackCounter + 1
	return c.SaveCounter()
}

func (c *Channel) ResetCounter() error {
	c.PlaybackCounter = 0
	return c.SaveCounter()
}

func (c Channel) SaveCounter() error {
	return client.Set(c.PlaybackCounterKey(), strconv.FormatInt(c.PlaybackCounter, 10), 0).Err()
}

func (c Channel) PlaybackCounterKey() string {
	return "broadcast-channel-" + c.Identifier + "-counter"
}
