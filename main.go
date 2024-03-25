package main

import (
	"memory-cache/cache"
	"time"
)

func main() {
	m := cache.NewMemoryCache("5GB")
	m.Set("1", 1, time.Minute)
	m.Set("2", 2, time.Minute)
	m.Set("3", 3, time.Minute)
	m.Set("4", 4, time.Minute)
	m.Set("5", 5, time.Minute)
}
