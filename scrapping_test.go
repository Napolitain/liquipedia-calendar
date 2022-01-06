package main

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
