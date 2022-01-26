package main

import (
	"testing"
)

func Test_byte2DToString1D(t *testing.T) {
	test := [][]byte{[]byte("test"), []byte("test2")}
	result := byte2DToString1D(test)
	if len(result) != 2 || result[0] != "test" || result[1] != "test2" {
		t.Fatal("Error while converting [][]byte to []string.")
	}
}

func Test_newQuery(t *testing.T) {
	query := newQuery("starcraft2", "maru")
	if query.game != "starcraft2" || len(query.players) != 1 || query.players[0] != "maru" {
		t.Fatal("New query is badly crafted: " + query.game + " " + query.players[0])
	}
}

func Test_newQueries(t *testing.T) {
	queries := newQueries("673d7374617263726166743226703d6d6172752c73657272616c3b673d6167656f66656d706972657326703d7468657669706572")
	if queries.data[0].game != "starcraft2" || len(queries.data[0].players) != 2 || queries.data[0].players[0] != "maru" || queries.data[0].players[1] != "serral" || queries.data[1].game != "ageofempires" || len(queries.data[1].players) != 1 || queries.data[1].players[0] != "theviper" {
		t.Fatal("New queries badly crafted.")
	}
}
