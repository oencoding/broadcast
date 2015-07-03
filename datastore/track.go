package datastore

// this file contains functions for dealing with tracks

import (
	"encoding/json"
	"github.com/omarqazi/broadcast/media"
)

func GetVideoTrack(trackId string) media.VideoTrack {
	trackType, err := client.Get(trackId + "-class").Result()
	if err != nil {
		trackType = "saved"
	}

	item := media.BlankTrack(trackType)
	trackData, err := client.Get("track-" + trackId + "-data").Result()
	if err != nil {
		return media.BlankTrack("")
	}

	item.Load(trackData)
	return item
}

func SaveVideoTrack(v media.VideoTrack) error {
	client.Set(v.Identifier()+"-class", v.Class(), 0)
	if jsonBytes, err := json.Marshal(v); err == nil {
		err := client.Set("track-"+v.Identifier()+"-data", string(jsonBytes), 0).Err()
		return err
	}

	return nil
}
