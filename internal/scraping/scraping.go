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

// getMatchesPage returns the appropriate matches page path for the given game
func getMatchesPage(game string) string {
	// League of Legends uses a different matches page
	if game == "leagueoflegends" {
		return MATCHES
	}
	// Default for all other games
	return UPCOMING_MATCHES
}

// GetFromLiquipedia function gets data from Liquipedia API and returns parsed HTML
func GetFromLiquipedia(game string) ([]byte, error) {
	matchesPage := getMatchesPage(game)
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
