package main

import (
	"encoding/json"
	"net/http"
)

// getFromLiquipedia function
func getFromLiquipedia(game string) (*http.Response, error) {
	url := BASE_URL + game + UPCOMING_MATCHES
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "liquipedia-calendar/developer (mxboucher@gmail.com)")
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
