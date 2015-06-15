package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"log"
	"net/http"
	"time"
)

var broadcastCursor = make(chan int)
var currentPlaylist string

type PlaylistGenerator struct {
	cursor chan int
}

func (pl *PlaylistGenerator) KeepPlaylistUpdated() {
	p, e := m3u8.NewMediaPlaylist(1000, 1000)
	if e != nil {
		log.Println("Error creating media playlist:", e)
		return
	}
	currentPlaylist = p.Encode().String()

	for {
		newMax := <-pl.cursor
		log.Println("new max", newMax)
		if err := p.Append(fmt.Sprintf("fileSequence%d.ts", newMax), 10.0, ""); err != nil {
			log.Println("Error appending item to playlist:", err, fmt.Sprintf("fileSequence%d.ts", newMax))
		}
		currentPlaylist = p.Encode().String()
	}
}

func (pl *PlaylistGenerator) Start() {
	pl.cursor = make(chan int, 1000)

	go pl.KeepPlaylistUpdated()
	for i := 0; i < 394; i++ {
		log.Println(i)
		pl.cursor <- (i % 5) + 30
		time.Sleep(10 * time.Second)
	}
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, currentPlaylist)
}
