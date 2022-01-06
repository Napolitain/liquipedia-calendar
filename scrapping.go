package main

import "net/http"

// getFromLiquipedia function
func getFromLiquipedia(game string) (*http.Response, error) {
	url := "https://liquipedia.net/" + game + "/api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches"
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
