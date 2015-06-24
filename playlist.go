package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"gopkg.in/redis.v1"
	"log"
	"net/http"
	"os/exec"
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
	generated := fmt.Sprintf("fileSequence%d.ts", seq)
	return generated
}

func (pl PlaylistGenerator) GeneratedVideoFileForSequence(seq int) string {
	prefix := ""
	pref := client.Get("broadcast-prefix").Val()
	prefix = pref

	generated := fmt.Sprintf("fileSequence%d.ts", seq)
	postProcess := fmt.Sprintf("fileSequence%d-post.ts", seq)
	sourceVideo := prefix + generated
	destVideo := prefix + postProcess

	currentTime := time.Now().Format("3:04 PM")

	twoClipsAgo := seq - 2
	if twoClipsAgo > 0 {
		mapKey := fmt.Sprintf("/fileSequence%d-post.ts", twoClipsAgo)
		log.Println("map key is", mapKey)
		if count, ok := lfs.Counter[mapKey]; ok {
			currentTime = fmt.Sprintf("%d active viewers", count)
		}
	}

	err := RenderTextToPNG(currentTime, "time.png")
	if err == nil {
		cmd := exec.Command("avconv", "-i", sourceVideo, "-vf", "movie=time.png [watermark];[in][watermark] overlay=0:0 [out]", "-y", "-map", "0", "-c:a", "copy", "-c:v", "mpeg2video", "-an", destVideo)
		err := cmd.Start()
		if err != nil {
			return sourceVideo
		}
		err = cmd.Wait()
		return destVideo
	}

	return sourceVideo
}

func (pl *PlaylistGenerator) KeepPlaylistUpdated() {
	p, e := m3u8.NewMediaPlaylist(1000, 1000)
	if e != nil {
		log.Println("Error creating media playlist:", e)
		return
	}
	currentPlaylist = p.Encode().String()

	for seqnum := 0; seqnum < 390; seqnum = <-pl.cursor {
		videoFile := pl.VideoFileForSequence(seqnum)
		if err := p.Append(videoFile, 10.0, ""); err != nil {
			log.Println("Error appending item to playlist:", err, fmt.Sprintf("fileSequence%d.ts", seqnum))
		}
		currentPlaylist = p.Encode().String()
	}
}

func (pl *PlaylistGenerator) Start() {
	pl.cursor = make(chan int, 1000)

	go pl.KeepPlaylistUpdated()
	for i := 1; i < 394; i++ {
		log.Println(i)
		pl.cursor <- i
		time.Sleep(10 * time.Second)
	}
}

func (pl PlaylistGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, currentPlaylist)
}
