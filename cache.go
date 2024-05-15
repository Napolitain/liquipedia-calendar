package main

import (
	"context"
	"google.golang.org/appengine/memcache"
	"time"
)

// saveToCache function is used to saveToCache data on Memcached server.
// A cached entry is a game entry valid for 1 hour. The whole entry is a HTML string to be filtered down later on for specific players.
// Reasonably fast, and extremely useful for reducing the load on the third party server Liquipedia.net.
// Likely will save a 100+ms per request.
func saveToCache(ctx context.Context, data string, game string) error {
	item := &memcache.Item{
		Key:        game,
		Value:      []byte(data),
		Expiration: time.Hour * 1,
	}
	return memcache.Set(ctx, item)
}

// getFromCache function is used to retrieve data from Memcached server. It returns HTML to be parsed and filtered down for specific players.
// Reasonably fast, and extremely useful for reducing the load on the third party server Liquipedia.net.
// Likely will save a 100+ms per request.
func getFromCache(ctx context.Context, game string) (*memcache.Item, error) {
	return memcache.Get(ctx, game)
}

// saveToCachePlayer function is used to saveToCache data on Memcached server. It is used to saveToCache data for a specific player for a specific game.
// The whole entry is a iCalendar string, ready to be sent to the user as a response directly.
// Super fast, very useful for reducing the load if the same player is queried multiple times (for example a superstar player).
func saveToCachePlayer(ctx context.Context, data string, queries string) error {
	item := &memcache.Item{
		Key:        queries,
		Value:      []byte(data),
		Expiration: time.Hour * 1,
	}
	return memcache.Set(ctx, item)
}

// getFromCachePlayer function is used to retrieve data from Memcached server. It returns iCalendar string ready to be sent to the user as a response directly.
// Super fast, very useful for reducing the load if the same player is queried multiple times (for example a superstar player).
func getFromCachePlayer(ctx context.Context, queries string) (*memcache.Item, error) {
	return memcache.Get(ctx, queries)
}

// getGamesFromCache function is used to retrieve the enabled games from cache
func getGamesFromCache(ctx context.Context) (*memcache.Item, error) {
	return memcache.Get(ctx, "!games")
}

// saveGamesToCache function is used to save the enabled games to cache
func saveGamesToCache(ctx context.Context, data string) error {
	item := &memcache.Item{
		Key:   "!games",
		Value: []byte(data),
	}
	return memcache.Set(ctx, item)
}
