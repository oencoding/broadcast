package main

import (
	"log"
	"net/http"
)

func main() {
	lfs := LogFileSystem{http.Dir(".")}
	fsHandler := http.FileServer(lfs)
	broadcast := PlaylistGenerator{}

	http.Handle("/", fsHandler)          // serve local directory on root
	http.Handle("/live.m3u8", broadcast) // serve generated playlist

	go broadcast.Start()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting HTTP server", err)
	}
}
