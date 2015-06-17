package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"gopkg.in/redis.v1"
	"log"
	"net/http"
	"time"
)

var broadcastCursor = make(chan int)
var currentPlaylist string
var client *redis.Client

func init() {
	client = redis.NewTCPClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := client.Ping().Result()
	log.Println(pong, err)
}

type PlaylistGenerator struct {
	cursor chan int
}

func (pl PlaylistGenerator) VideoFileForSequence(seq int) string {
	prefix := ""
	pref := client.Get("broadcast-prefix").Val()
	prefix = pref

	generated := fmt.Sprintf("fileSequence%d.ts", seq)
	return prefix + generated
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
		if err := p.Append(pl.VideoFileForSequence(newMax), 10.0, ""); err != nil {
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
		pl.cursor <- i
		time.Sleep(10 * time.Second)
	}
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, currentPlaylist)
}
