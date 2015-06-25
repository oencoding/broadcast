package datastore

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

func (c Channel) PlaybackCounterKey() string {
	return "broadcast-channel-" + c.Identifier + "-counter"
}
