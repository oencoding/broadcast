package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/omarqazi/broadcast/datastore"
	"log"
	"net/http"
	"strings"
	"time"
)

var broadcastCursor = make(chan int)
var currentPlaylist string
var allChannels map[string]*datastore.Channel

type PlaylistGenerator struct {
	cursor chan int
}

func (pl PlaylistGenerator) VideoFileForSequence(seq int) string {
	generated := fmt.Sprintf("http://www.smick.tv/media/truedetectives2e1movie%05d.ts", seq)
	return generated
}

func (pl *PlaylistGenerator) KeepPlaylistUpdated() {
	p, e := m3u8.NewMediaPlaylist(1000, 1000)
	if e != nil {
		log.Println("Error creating media playlist:", e)
		return
	}
	currentPlaylist = p.Encode().String()

	for seqnum := 1; seqnum < 1854; seqnum = <-pl.cursor {
		videoFile := pl.VideoFileForSequence(seqnum)
		if err := p.Append(videoFile, 5.0, ""); err != nil {
			log.Println("Error appending item to playlist:", err, fmt.Sprintf("movie2m%5d.ts", seqnum))
		}
		currentPlaylist = p.Encode().String()
	}
}

func (pl *PlaylistGenerator) Start() {
	pl.cursor = make(chan int, 1000)

	go pl.KeepPlaylistUpdated()
	for i := 1; i < 1854; i++ {
		log.Println(i)
		pl.cursor <- i
		time.Sleep(5 * time.Second)
	}
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelId := strings.TrimSuffix(r.URL.Path, ".m3u8")
	channel, ok := allChannels[channelId]
	if !ok { // If this is the first time the channel is requested
		channel = datastore.GetChannel(channelId)
		go channel.AdvanceEvery(5 * time.Second)
		allChannels[channelId] = &channel
	}
	channel := datastore.GetChannel(cha)
	fmt.Fprintln(w, currentPlaylist)
}
