package main

import (
	"log"
	"net/http"
)

var lfs = LogFileSystem{http.Dir("."), make(map[string]int)}
const channelRoute = "/channel/"

func main() {
	fsHandler := http.FileServer(lfs)
	broadcast := PlaylistGenerator{}
	channelHandler := http.StripPrefix(channelRoute, broadcast)

	http.Handle("/", fsHandler) // serve current directory on root
	http.Handle(channelRoute, channelHandler)

	go broadcast.Start()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting HTTP server", err)
	}
}
