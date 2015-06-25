package main

import (
	"log"
	"net/http"
)

const channelRoute = "/channel/"

var lfs = LogFileSystem{http.Dir("."), make(map[string]int)}

// when the program starts:
// 1. Setup up routes to:
//    a) Serve the current directory
//    b) Serve playlists off the channel route
// 2. Start the broadcast timer
// 3. Start a HTTP server on port 8080

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
