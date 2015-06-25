package datastore

import (
	"testing"
)

func TestPing(t *testing.T) {
	if _, err := client.Ping().Result(); err != nil {
		t.Error("Error pinging redis:", err)
	}
}
