package main

import (
	"log"
	"net/http"
)

func main() {
	lfs := LogFileSystem{http.Dir(".")}
	fsHandler := http.FileServer(lfs)
	http.Handle("/", fsHandler)                    // serve local directory on root
	http.Handle("/live.m3u8", PlaylistGenerator{}) // serve generated playlist

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting HTTP server", err)
	}
}
