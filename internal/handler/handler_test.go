package handler

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
	if query.Game != "starcraft2" || len(query.Players) != 1 || query.Players[0] != "maru" {
		t.Fatal("New query is badly crafted: " + query.Game + " " + query.Players[0])
	}
}

func Test_NewQueries(t *testing.T) {
	queries, err := NewQueries("673d7374617263726166743226703d6d6172752c73657272616c3b673d6167656f66656d706972657326703d7468657669706572")
	if err != nil {
		t.Fatal("Error raised.")
	}
	if queries.Data[0].Game != "starcraft2" || len(queries.Data[0].Players) != 2 || queries.Data[0].Players[0] != "maru" || queries.Data[0].Players[1] != "serral" || queries.Data[1].Game != "ageofempires" || len(queries.Data[1].Players) != 1 || queries.Data[1].Players[0] != "theviper" {
		t.Fatal("New queries badly crafted.")
	}
}

func Test_NewQueries_1(t *testing.T) {
	queries, err := NewQueries("673d7374617263726166743226703d4d6172752c53657272616c")
	if err != nil {
		t.Fatal("Error raised.")
	}
	if queries.Data[0].Game != "starcraft2" || len(queries.Data[0].Players) != 2 || queries.Data[0].Players[0] != "Maru" || queries.Data[0].Players[1] != "Serral" {
		t.Fatal("New queries badly crafted.")
	}
}

// Test if error is raised when hexadecimal string is badly crafted.
func Test_NewQueries_badHex(t *testing.T) {
	_, err := NewQueries("673d737461263726166743226703d4d6172752c53657272616c")
	if err == nil {
		t.Fatal("Error not raised.")
	}
}
