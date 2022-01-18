package pkg

import (
	"io/ioutil"
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log(string(body[:]))
}

func TestParseJSON(t *testing.T) {
	input := []byte(`{"parse":{"title":"Liquipedia:Upcoming and ongoing matches","pageid":64908,"revid":2145797,"text":{"*":"<div class=\"mw-parser-output\">"}}}`)
	bytes, err := parseJSON(input)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log(bytes)
}
