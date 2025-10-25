package handler

import (
	"bytes"
	"encoding/hex"
	"log/slog"
)

var logger = slog.Default()

// Convert [][]byte to []string
func byte2DToString1D(byte2d [][]byte) []string {
	var p []string
	for i := 0; i < len(byte2d); i++ {
		p = append(p, string(byte2d[i][:]))
	}
	return p
}

// Query struct
// ex: g=starcraft2&p=maru,serral
type Query struct {
	Game    string
	Players []string
}

// Create a Query struct
func newQuery(game string, players string) *Query {
	decodeString := bytes.Split([]byte(players), []byte(","))
	p := byte2DToString1D(decodeString)
	return &Query{Game: game, Players: p}
}

// Queries struct
// ex: Queries=hexadecimal (g=starcraft2&p=maru,serral;g=ageofempires&p=theviper)
type Queries struct {
	Data []Query
}

// NewQueries creates a Queries struct (made of multiple Query)
func NewQueries(query string) (*Queries, error) {
	decodeString, err := hex.DecodeString(query)
	if err != nil {
		logger.Error("Failed to decode query hex", "error", err)
		return nil, err
	}
	queries := bytes.Split(decodeString, []byte(";")) // g=starcraft2&p=maru,serral
	var result Queries
	for i := 0; i < len(queries); i++ {
		q := bytes.Split(queries[i], []byte("&")) // g=starcraft2
		game := bytes.Split(q[0], []byte("="))
		players := bytes.Split(q[1], []byte("="))
		result.Data = append(result.Data, *newQuery(string(game[1][:]), string(players[1][:])))
	}
	return &result, nil
}
