package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      []byte
	expiration time.Time
}

type Cache struct {
	items map[string]*cacheItem
	mu    sync.RWMutex
}

var globalCache *Cache

// Init initializes the global cache
func Init() {
	globalCache = &Cache{
		items: make(map[string]*cacheItem),
	}
	// Start a goroutine to clean up expired items periodically
	go globalCache.cleanup()
}

// cleanup removes expired items from the cache
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// set adds an item to the cache with an expiration time
func (c *Cache) set(key string, value []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(duration),
	}
}

// get retrieves an item from the cache
func (c *Cache) get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	if time.Now().After(item.expiration) {
		return nil, false
	}
	return item.value, true
}

// SetGameData saves game data to cache for 1 hour
func SetGameData(game string, data []byte) {
	globalCache.set(game, data, time.Hour)
}

// GetGameData retrieves game data from cache
func GetGameData(game string) []byte {
	data, exists := globalCache.get(game)
	if !exists {
		return nil
	}
	return data
}

// SetPlayerCalendar saves player calendar to cache for 1 hour
func SetPlayerCalendar(queries string, data string) {
	globalCache.set(queries, []byte(data), time.Hour)
}

// GetPlayerCalendar retrieves player calendar from cache
func GetPlayerCalendar(queries string) (string, error) {
	data, exists := globalCache.get(queries)
	if !exists {
		return "", nil
	}
	return string(data), nil
}

// SetGames saves the list of games to cache (no expiration)
func SetGames(data string) {
	globalCache.set("!games", []byte(data), 24*time.Hour)
}

// GetGames retrieves the list of games from cache
func GetGames() string {
	data, exists := globalCache.get("!games")
	if !exists {
		return ""
	}
	return string(data)
}
