package pkg

import (
	"context"
	"io/ioutil"
)

// getData returns data in []byte format from either cache or scrapping
func getData(ctx context.Context, game string) ([]byte, error) {
	// Get data from cache server
	item, err := getFromCache(ctx, game)
	if err != nil {
		// If fail, get data from scrapping
		response, err := getFromLiquipedia(game)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, err
		}

		// Convert from io to []byte
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		// parse JSON
		body, err = parseJSON(body)
		if err != nil {
			return nil, err
		}

		// Save to cache server
		err = saveToCache(ctx, string(body[:]), game)
		if err != nil {
			return nil, err
		}
		return body, nil
	} else {
		return item.Value, nil
	}
}
