//package configuration provides functions for determing the proper
//configuration of the broadcast server
package configuration

import (
	"os"
)

// constant RedisEnvironmentVariable sets the name of the enviornment variable
// used to get the redis address
const RedisEnvironmentVariable = "REDIS_ADDR"

// function RedisServerAddress returns the address of the configured redis server
// It gets its value from the environment variable REDIS_ADDR, defaults to localhost:6379
func RedisServerAddress() (addr string) {
	addr = "localhost:6379" // default value
	if envVar := os.Getenv(RedisEnvironmentVariable); envVar != "" {
		addr = envVar
	}
	return
}
