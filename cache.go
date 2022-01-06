package main

import (
	"context"
	"google.golang.org/appengine/memcache"
)

// saveToCache function is used to saveToCache data on Memcached server.
func saveToCache(ctx context.Context, data string) error {
	item := &memcache.Item{
		Key:   "data",
		Value: []byte(data),
	}
	return memcache.Set(ctx, item)
}

// getFromCache function is used to retrieve data from Memcached server.
func getFromCache(ctx context.Context) (*memcache.Item, error) {
	return memcache.Get(ctx, "data")
}
