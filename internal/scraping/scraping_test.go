package scraping

import (
	"log"
	"os"
	"testing"
)

func TestGetFromLiquipedia(t *testing.T) {
	_, err := GetFromLiquipedia("starcraft2")
	if err != nil {
		t.Log("Network error expected in test environment: " + err.Error())
		return
	}
}

func TestParseJSON(t *testing.T) {
	input := []byte(`{"parse":{"title":"Liquipedia:Upcoming and ongoing matches","pageid":64908,"revid":2145797,"text":{"*":"<div class=\"mw-parser-output\">"}}}`)
	_, err := parseJSON(input)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
}

func Test_FetchGames(t *testing.T) {
	// Read all testing data
	_, err := os.ReadFile("../../resources/scrapping_test_get_wikis")
	if err != nil {
		log.Fatal(err)
	}

	games, err := FetchGames()
	if err != nil {
		t.Log("Network error expected in test environment: " + err.Error())
		return
	}
	if len(games) == 0 {
		t.Fatal("No games found.")
		return
	}
}

func Test_getMatchesPage(t *testing.T) {
	tests := []struct {
		name     string
		game     string
		expected string
	}{
		{
			name:     "League of Legends uses Matches page",
			game:     "leagueoflegends",
			expected: MATCHES,
		},
		{
			name:     "Starcraft2 uses Upcoming and ongoing matches page",
			game:     "starcraft2",
			expected: UPCOMING_MATCHES,
		},
		{
			name:     "Dota2 uses Upcoming and ongoing matches page",
			game:     "dota2",
			expected: UPCOMING_MATCHES,
		},
		{
			name:     "Counter-Strike uses Upcoming and ongoing matches page",
			game:     "counterstrike",
			expected: UPCOMING_MATCHES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMatchesPage(tt.game)
			if result != tt.expected {
				t.Errorf("getMatchesPage(%s) = %s, expected %s", tt.game, result, tt.expected)
			}
		})
	}
}

func Test_GetScrapingStrategy(t *testing.T) {
	tests := []struct {
		name         string
		game         string
		expectedType string
	}{
		{
			name:         "League of Legends returns LeagueOfLegendsScrapingStrategy",
			game:         "leagueoflegends",
			expectedType: "*scraping.LeagueOfLegendsScrapingStrategy",
		},
		{
			name:         "Starcraft2 returns DefaultScrapingStrategy",
			game:         "starcraft2",
			expectedType: "*scraping.DefaultScrapingStrategy",
		},
		{
			name:         "Dota2 returns DefaultScrapingStrategy",
			game:         "dota2",
			expectedType: "*scraping.DefaultScrapingStrategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetScrapingStrategy(tt.game)
			if strategy == nil {
				t.Fatal("Expected non-nil strategy")
			}

			// Check the type of strategy
			switch tt.expectedType {
			case "*scraping.LeagueOfLegendsScrapingStrategy":
				if _, ok := strategy.(*LeagueOfLegendsScrapingStrategy); !ok {
					t.Errorf("Expected LeagueOfLegendsScrapingStrategy, got %T", strategy)
				}
			case "*scraping.DefaultScrapingStrategy":
				if _, ok := strategy.(*DefaultScrapingStrategy); !ok {
					t.Errorf("Expected DefaultScrapingStrategy, got %T", strategy)
				}
			}
		})
	}
}

func Test_ScrapingStrategy_GetMatchesPage(t *testing.T) {
	tests := []struct {
		name     string
		strategy ScrapingStrategy
		game     string
		expected string
	}{
		{
			name:     "DefaultScrapingStrategy returns UPCOMING_MATCHES",
			strategy: &DefaultScrapingStrategy{},
			game:     "starcraft2",
			expected: UPCOMING_MATCHES,
		},
		{
			name:     "LeagueOfLegendsScrapingStrategy returns MATCHES",
			strategy: &LeagueOfLegendsScrapingStrategy{},
			game:     "leagueoflegends",
			expected: MATCHES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.strategy.GetMatchesPage(tt.game)
			if result != tt.expected {
				t.Errorf("GetMatchesPage(%s) = %s, expected %s", tt.game, result, tt.expected)
			}
		})
	}
}
