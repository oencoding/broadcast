package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"log"
	"net/http"
)

type PlaylistGenerator struct {
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, e := m3u8.NewMediaPlaylist(100, 200)
	if e != nil {
		log.Println("Error creating media playlist:", e)
	}

	for i := 0; i < 394; i++ {
		if e = p.Append(fmt.Sprintf("fileSequence%d.ts", i), 10.0, ""); e != nil {
			log.Println("Error appending item to playlist:", e)
		}
	}

	fmt.Fprintln(w, p.Encode().String())
}
