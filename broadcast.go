package main

import (
	"log"
	"net/http"
)

var lfs = LogFileSystem{http.Dir("."), make(map[string]int)}

func main() {
	fsHandler := http.FileServer(lfs)
	broadcast := PlaylistGenerator{}

	http.Handle("/", fsHandler)          // serve local directory on root
	http.Handle("/live.m3u8", broadcast) // serve generated playlist

	go broadcast.Start()

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Println("Error starting HTTP server", err)
	}
}
