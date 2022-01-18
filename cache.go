package main

import (
	"context"
	"google.golang.org/appengine/memcache"
	"time"
)

// saveToCache function is used to saveToCache data on Memcached server.
func saveToCache(ctx context.Context, data string, game string) error {
	item := &memcache.Item{
		Key:        game,
		Value:      []byte(data),
		Expiration: time.Hour * 3,
	}
	return memcache.Set(ctx, item)
}

// getFromCache function is used to retrieve data from Memcached server.
func getFromCache(ctx context.Context, game string) (*memcache.Item, error) {
	return memcache.Get(ctx, game)
}
