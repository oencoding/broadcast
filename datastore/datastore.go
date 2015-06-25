// package datastore provides models for data structures,
// and functions for getting and saving them to Redis
package datastore

import (
	"github.com/omarqazi/broadcast/configuration"
	"gopkg.in/redis.v3"
	"log"
)

var client *redis.Client

func init() {
	options := &redis.Options{Addr: configuration.RedisServerAddress()}
	client = redis.NewClient(options)
	if _, err := client.Ping().Result(); err != nil {
		log.Fatalln("Error connecting to redis:", err)
	}
}
