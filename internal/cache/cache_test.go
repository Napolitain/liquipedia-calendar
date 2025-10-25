package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	// Initialize cache
	Init()

	// Test SetGameData and GetGameData
	testData := []byte("test data")
	SetGameData("testgame", testData)
	
	result := GetGameData("testgame")
	if result == nil {
		t.Fatal("Failed to retrieve cached game data")
	}
	if string(result) != string(testData) {
		t.Fatal("Retrieved data doesn't match stored data")
	}

	// Test SetPlayerCalendar and GetPlayerCalendar
	testCalendar := "test calendar"
	SetPlayerCalendar("testquery", testCalendar)
	
	calendar, err := GetPlayerCalendar("testquery")
	if err != nil {
		t.Fatal("Error retrieving player calendar:", err)
	}
	if calendar != testCalendar {
		t.Fatal("Retrieved calendar doesn't match stored calendar")
	}

	// Test SetGames and GetGames
	testGames := "game1,game2,game3"
	SetGames(testGames)
	
	games := GetGames()
	if games != testGames {
		t.Fatal("Retrieved games don't match stored games")
	}

	// Test cache expiration
	globalCache.set("tempkey", []byte("temp value"), 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	
	expired, exists := globalCache.get("tempkey")
	if exists {
		t.Fatal("Expired cache item should not be retrievable")
	}
	if expired != nil {
		t.Fatal("Expired cache item should return nil")
	}

	// Test cache miss
	missing := GetGameData("nonexistent")
	if missing != nil {
		t.Fatal("Non-existent cache item should return nil")
	}
}
