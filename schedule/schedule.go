// until i have a better UI built, schedule helps schedule programming into redis
package main

import (
	"github.com/omarqazi/broadcast/datastore"
	"github.com/omarqazi/broadcast/media"
	"log"
)

func main() {
	channel, _ := datastore.GetChannel("live")
	smicktv := &media.HLSTrack{
		Id:              "meerkat",
		PlaylistURL:     "http://cdn.meerkatapp.co/broadcast/c42beb6b-4a19-4227-8cbd-303a0e796d28/live.m3u8",
		PlaybackCounter: 0,
	}

	datastore.SaveVideoTrack(smicktv)
	pl := datastore.PlaylistItem{
		TrackId: "meerkat",
		Loop:    false,
	}

	if err := channel.PushItem(pl); err != nil {
		log.Println("Error:", err)
	}
}
