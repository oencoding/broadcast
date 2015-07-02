// until i have a better UI built, schedule helps schedule programming into redis
package main

import "github.com/omarqazi/broadcast/datastore"

func main() {
	channel, _ := datastore.GetChannel("live")
	phil := &datastore.PlaylistItem{
		URLFormat: "http://www.smick.tv/media/away/fileSequence%d.ts",
		StartAt: 0,
		EndAt: 53,
		Loop: true,
	}
	
	channel.PushItem(phil)
}

func xmain() {
    channel, _ := datastore.GetChannel("live")
    ad := &datastore.PlaylistItem{
        URLFormat: "http://www.smick.tv/media/smick/fileSequence%d.ts",
        StartAt: 0,
        EndAt: 36,
        Loop: false,
    }
    
    episodeBreaks := map[string][]int64{
        "dexters1e2" : []int64{87,170,258},
        "dexters1e7" : []int64{88,172,259},
        "dexters1e10" : []int64{87,170,258},
    }
    
    scheduleEpisode := func(episodeName string) {
        breaks := episodeBreaks[episodeName]
        lastend := int64(0)
        for i := range breaks {
            endAt := breaks[i] - 1
            channel.PushItem(mediaItem(episodeName,lastend,endAt))
            lastend = breaks[i]
            channel.PushItem(ad)
        }
    }
    
    for j := 0;j < 100;j++ {
        for k, _ := range episodeBreaks {
            scheduleEpisode(k)
        }   
    } 
}

func mediaItem(folder string,start int64,end int64) *datastore.PlaylistItem {
    return &datastore.PlaylistItem{
        URLFormat: "http://www.smick.tv/media/" + folder + "/fileSequence%d.ts",
        StartAt: start,
        EndAt: end,
        Loop: false,
    }
}