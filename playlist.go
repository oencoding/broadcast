package main

import (
	"fmt"
	"github.com/omarqazi/broadcast/datastore"
	"net/http"
	"strings"
	"time"
)

var allChannels map[string]*datastore.Channel = make(map[string]*datastore.Channel)

type PlaylistGenerator struct {
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	checkCORSHeader(w, r)
	var err error
	channelId := strings.TrimSuffix(r.URL.Path, ".m3u8")
	channel, ok := allChannels[channelId]
	if !ok { // If this is the first time the channel is requested
		channel, err = datastore.GetChannel(channelId)
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}

		go channel.AdvanceEvery(5*time.Second, nil)
		allChannels[channelId] = channel
	}

	fmt.Fprintln(w, channel.PlaylistData())
}

func checkCORSHeader(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	url, err := url.Parse(origin)
	if err != nil {
		errorOut(w, r, "Invalid Origin header")
	}
	host := url.Host
	w.Header().Add("Access-Control-Allow-Origin", origin)
}
