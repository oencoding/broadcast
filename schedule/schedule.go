// until i have a better UI built, schedule helps schedule programming into redis
package main

import (
	"github.com/omarqazi/broadcast/datastore"
	"github.com/omarqazi/broadcast/media"
	"log"
)

func main() {
	channel, _ := datastore.GetChannel("live")
	prydz := media.SavedTrack{
		Id:             "prydz",
		URLFormat:      "http://www.smick.tv/media/prydz/fileSequence%d.ts",
		StartAt:        0,
		EndAt:          157,
		Loop:           true,
		TargetDuration: 5.0,
	}

	if err := channel.PushItem(prydz); err != nil {
		log.Println("Error:", err)
	}
}

func mediaItem(folder string, start int64, end int64) *datastore.PlaylistItem {
	return &datastore.PlaylistItem{
		URLFormat: "http://www.smick.tv/media/" + folder + "/fileSequence%d.ts",
		StartAt:   start,
		EndAt:     end,
		Loop:      false,
	}
}
