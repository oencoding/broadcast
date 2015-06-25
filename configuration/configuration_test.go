package configuration

import (
    "testing"
    "os"
)

func TestRedisServerAddress(t *testing.T) {
    oldEnv := os.Getenv(RedisEnvironmentVariable)
    os.Setenv(RedisEnvironmentVariable,"")
    
    if rv := RedisServerAddress(); rv != "localhost:6379" {
        t.Error("Error: Expected default redis server address of localhost:6379 but got",rv)
    }
    
    const someValue = "192.168.1.2:6379"
    os.Setenv(RedisEnvironmentVariable,someValue)
    if rv := RedisServerAddress(); rv != someValue {
        t.Error("Error: expected redis server address of",someValue,"but got",rv)
    }
    
    os.Setenv(RedisEnvironmentVariable,oldEnv)
}