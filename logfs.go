package main

// logfs.go just implements a net/http FileSystem that logs requests for files
// this may help with debugging

import (
	"log"
	"net/http"
)

// LogFileSystem wraps http.FileSystem, but logs requests
type LogFileSystem struct {
	fs http.FileSystem
}

func (l LogFileSystem) Open(name string) (f http.File, e error) {
	log.Println(name) // do our magic

	// then, do whatever http.FileSystem does
	f, e = l.fs.Open(name)
	return
}
