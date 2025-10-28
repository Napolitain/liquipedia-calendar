package scraping

import (
	"encoding/json"
	"log/slog"

	"github.com/valyala/fasthttp"
	"k8s.io/apimachinery/pkg/util/sets"
)

const BASE_URL = "https://liquipedia.net/"
const UPCOMING_MATCHES = "/api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches"
const MATCHES = "/api.php?action=parse&format=json&page=Liquipedia:Matches"

var logger = slog.Default()

// ScrapingStrategy defines the interface for different scraping strategies.
// Implement this interface to create custom scraping strategies for games that
// require different API endpoints or scraping logic.
//
// Example usage:
//
//	type CustomGameStrategy struct{}
//
//	func (s *CustomGameStrategy) GetMatchesPage(game string) string {
//	    return "/custom/api/endpoint"
//	}
//
//	func (s *CustomGameStrategy) ScrapeData(game string) ([]byte, error) {
//	    // Custom implementation
//	    matchesPage := s.GetMatchesPage(game)
//	    return fetchFromLiquipedia(game, matchesPage)
//	}
//
// Then update GetScrapingStrategy to return your custom strategy:
//
//	if game == "customgame" {
//	    return &CustomGameStrategy{}
//	}
type ScrapingStrategy interface {
	// GetMatchesPage returns the API endpoint path for the given game
	GetMatchesPage(game string) string
	// ScrapeData fetches and returns the HTML data from Liquipedia
	ScrapeData(game string) ([]byte, error)
}

// DefaultScrapingStrategy implements the default scraping strategy for most games
type DefaultScrapingStrategy struct{}

// GetMatchesPage returns the default matches page path
func (s *DefaultScrapingStrategy) GetMatchesPage(game string) string {
	return UPCOMING_MATCHES
}

// ScrapeData fetches data from Liquipedia using the default strategy
func (s *DefaultScrapingStrategy) ScrapeData(game string) ([]byte, error) {
	matchesPage := s.GetMatchesPage(game)
	return fetchFromLiquipedia(game, matchesPage)
}

// LeagueOfLegendsScrapingStrategy implements the scraping strategy for League of Legends
type LeagueOfLegendsScrapingStrategy struct{}

// GetMatchesPage returns the League of Legends specific matches page path
func (s *LeagueOfLegendsScrapingStrategy) GetMatchesPage(game string) string {
	return MATCHES
}

// ScrapeData fetches data from Liquipedia using the League of Legends strategy
func (s *LeagueOfLegendsScrapingStrategy) ScrapeData(game string) ([]byte, error) {
	matchesPage := s.GetMatchesPage(game)
	return fetchFromLiquipedia(game, matchesPage)
}

// GetScrapingStrategy returns the appropriate scraping strategy for the given game
func GetScrapingStrategy(game string) ScrapingStrategy {
	// League of Legends uses a different strategy
	if game == "leagueoflegends" {
		return &LeagueOfLegendsScrapingStrategy{}
	}
	// Default strategy for all other games
	return &DefaultScrapingStrategy{}
}

// getMatchesPage returns the appropriate matches page path for the given game
// Deprecated: Use GetScrapingStrategy(game).GetMatchesPage(game) instead
func getMatchesPage(game string) string {
	strategy := GetScrapingStrategy(game)
	return strategy.GetMatchesPage(game)
}

// GetFromLiquipedia function gets data from Liquipedia API and returns parsed HTML
// It uses the appropriate scraping strategy based on the game
func GetFromLiquipedia(game string) ([]byte, error) {
	strategy := GetScrapingStrategy(game)
	return strategy.ScrapeData(game)
}

// fetchFromLiquipedia is a helper function that performs the actual HTTP request and parsing
func fetchFromLiquipedia(game string, matchesPage string) ([]byte, error) {
	url := BASE_URL + game + matchesPage
	logger.Info("GET request to Liquipedia", "url", url)
	
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("User-Agent", "liquipedia-calendar (https://github.com/Napolitain/liquipedia-calendar mxboucher@gmail.com)")

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		logger.Error("Error while getting data from Liquipedia", "status", resp.StatusCode())
		return nil, err
	}

	body := resp.Body()

	// parse JSON
	data, err := parseJSON(body)
	if err != nil {
		return nil, err
	}
	
	logger.Info("Retrieved data from Liquipedia", "size", len(data))
	return data, nil
}

// parseJSON parses JSON response from liquipedia and returns []byte containing HTML.
func parseJSON(in []byte) ([]byte, error) {
	// Declared an empty map interface
	var result map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal(in, &result)
	if err != nil {
		return nil, err
	}
	parse := result["parse"].(map[string]interface{})["text"].(map[string]interface{})["*"].(string)
	return []byte(parse), nil
}

// FetchGames retrieves the list of games supported by Liquipedia API
func FetchGames() (sets.String, error) {
	url := "https://liquipedia.net/api.php?action=listwikis"
	
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("User-Agent", "liquipedia-calendar (https://github.com/Napolitain/liquipedia-calendar mxboucher@gmail.com)")

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		logger.Error("Failed to fetch games", "error", err)
		return nil, err
	}

	body := resp.Body()
	
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Failed to parse games JSON", "error", err)
		return nil, err
	}
	
	// Parse JSON
	wikis := result["allwikis"].(map[string]interface{})
	// Create a set of games to avoid duplicates
	games := sets.NewString()
	// Get the ID of each game (each game, not more, is allowed). Every game ID is the key in the allwikis map.
	for key := range wikis {
		games.Insert(key)
	}
	return games, nil
}
