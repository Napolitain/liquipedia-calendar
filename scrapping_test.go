package main

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestGetFromLiquipedia(t *testing.T) {
	response, err := getFromLiquipedia("starcraft2")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	if response.StatusCode != 200 {
		t.Fatal(response.Status)
		return
	}
	_, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
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

func Test_fetchGames(t *testing.T) {
	// Read all testing data
	_, err := os.ReadFile("resources/scrapping_test_get_wikis")
	if err != nil {
		log.Fatal(err)
	}

	games, err := fetchGames()
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	if len(games) == 0 {
		t.Fatal("No games found.")
		return
	}
}
