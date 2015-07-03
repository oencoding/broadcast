// until i have a better UI built, schedule helps schedule programming into redis
package main

import (
	"github.com/omarqazi/broadcast/datastore"
	"github.com/omarqazi/broadcast/media"
	"log"
)

func main() {
	channel, _ := datastore.GetChannel("live")
	prydz := &media.SavedTrack{
		Id:             "phil",
		URLFormat:      "http://www.smick.tv/media/philshow/fileSequence%d.ts",
		StartAt:        0,
		EndAt:          14,
		TargetDuration: 5.0,
	}

	datastore.SaveVideoTrack(prydz)
	pl := datastore.PlaylistItem{
		TrackId: "phil",
		Loop:    true,
	}

	if err := channel.PushItem(pl); err != nil {
		log.Println("Error:", err)
	}
}
