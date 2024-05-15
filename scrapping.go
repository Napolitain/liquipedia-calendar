package main

import (
	"encoding/json"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"
)

// Guidelines for Liquipedia API usage
// https://liquipedia.net/commons/Liquipedia:API_Usage_Guidelines

// getFromLiquipedia function
func getFromLiquipedia(game string) (*http.Response, error) {
	url := BASE_URL + game + UPCOMING_MATCHES
	logger.Println("GET " + url)
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "liquipedia-calendar (https://github.com/Napolitain/liquipedia-calendar mxboucher@gmail.com)")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
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

// List of games supported by Liquipedia API : important to avoid not only errors, but attacks.
func fetchGames() (sets.String, error) {
	// Prepare GET request to https://liquipedia.net/api.php?action=listwikis
	request, err := http.NewRequest("GET", "https://liquipedia.net/api.php?action=listwikis", nil)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Set User-Agent to avoid being blocked by Liquipedia
	request.Header.Set("User-Agent", "liquipedia-calendar (https://github.com/Napolitain/liquipedia-calendar mxboucher@gmail.com)")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Read the body of the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Parse JSON
	wikis := result["allwikis"].(map[string]interface{})
	// Create a set of games to avoid duplicates
	games := sets.NewString()
	// Get the ID of each game (each game, not more, is allowed). Every game ID is the key in the allwikis map.
	for i := 0; i < len(wikis); i++ {
		// Iterate over the map wikis
		for key, _ := range wikis {
			games.Insert(key)
		}
	}
	return games, nil
}
