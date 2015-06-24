package main

// logfs.go just implements a net/http FileSystem that logs requests for files
// this may help with debugging

import (
	"log"
	"net/http"
)

// LogFileSystem wraps http.FileSystem, but logs requests
type LogFileSystem struct {
	fs      http.FileSystem
	Counter map[string]int
}

func (l LogFileSystem) Open(name string) (f http.File, e error) {
	log.Println(name) // do our magic

	if val, ok := l.Counter[name]; ok {
		l.Counter[name] = val + 1
	} else {
		l.Counter[name] = 1
	}

	log.Println(l.Counter)

	// then, do whatever http.FileSystem does
	f, e = l.fs.Open(name)
	return
}
