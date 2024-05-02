package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"io"
	"strconv"
	"time"
)

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
	game    string
	players []string
}

// Create a Query struct
func newQuery(game string, players string) *Query {
	decodeString := bytes.Split([]byte(players), []byte(","))
	p := byte2DToString1D(decodeString)
	return &Query{game: game, players: p}
}

// Queries struct
// ex: Queries=hexadecimal (g=starcraft2&p=maru,serral;g=ageofempires&p=theviper)
type Queries struct {
	data []Query
}

type Calendar interface {
	createCalendar(document *goquery.Document, player Query) (*ics.Calendar, error)
}

// Create a Queries struct (made of multiple Query)
func newQueries(query string) *Queries {
	decodeString, err := hex.DecodeString(query)
	if err != nil {
		return nil
	}
	queries := bytes.Split(decodeString, []byte(";")) // g=starcraft2&p=maru,serral
	var result Queries
	for i := 0; i < len(queries); i++ {
		q := bytes.Split(queries[i], []byte("&")) // g=starcraft2
		game := bytes.Split(q[0], []byte("="))
		players := bytes.Split(q[1], []byte("="))
		result.data = append(result.data, *newQuery(string(game[1][:]), string(players[1][:])))
	}
	return &result
}

// getData returns data in []byte format from either cache or scrapping
func getData(ctx context.Context, game string) ([]byte, bool, error) {
	// Get data from cache server
	item, err := getFromCache(ctx, game)
	if err != nil {
		// If fail, get data from scrapping
		response, err := getFromLiquipedia(game)
		if err != nil {
			logger.Println(err)
			return nil, false, err
		}
		if response.StatusCode != 200 {
			logger.Println("Error while getting data from Liquipedia. Code " + response.Status)
			return nil, false, err
		}

		// Convert from io to []byte
		body, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Println(err)
			return nil, false, err
		}

		// parse JSON
		body, err = parseJSON(body)
		logger.Println("Length of body: " + strconv.Itoa(len(body)))
		if err != nil {
			logger.Println(err)
			return nil, false, err
		}

		// Save to cache server
		err = saveToCache(ctx, string(body[:]), game)
		if err != nil {
			logger.Println(err)
			return nil, false, err
		}
		return body, false, nil
	} else {
		return item.Value, true, nil
	}
}

/**
 * createCalendar inherits from Queries struct
 */
func (queries Queries) createCalendar(document *goquery.Document, player Query) (*ics.Calendar, error) {
	// Create iCalendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//liquipedia-calendar//en")
	cal.SetVersion("2.0")
	cal.SetLastModified(time.Now())
	// Create events
	matches := document.Find(".infobox_matches_content")
	var UIDs []string
	for i := 0; i < matches.Size(); i++ {
		// Get event info
		teamleft := matches.Eq(i).Find(".team-left a").Eq(0)
		teamright := matches.Eq(i).Find(".team-right span:not(.flag):not(.team-template-image):not(.team-template-team-short) a").Eq(0)
		teamleft_text := teamleft.Text() // TODO: test all match cases
		teamright_text := teamright.Text()
		teamleft_title := teamleft.AttrOr("title", "")
		teamright_title := teamright.AttrOr("title", "")
		contains := false
		for _, p := range player.players {
			// Filter out events that don't contain at least one player in players
			if teamleft_text == p || teamright_text == p || teamleft_title == p || teamright_title == p {
				contains = true
				break
			}
		}
		if contains == false {
			continue
		}
		matchFormat := matches.Eq(i).Find(".versus abbr").Eq(0).Text()
		timestampStr, exist := matches.Eq(i).Find("[data-timestamp]").Eq(0).Attr("data-timestamp")
		if exist != true {
			return nil, errors.New("Timestamp doesn't exist.")
		}
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return nil, errors.New("Error while converting string to int.")
		}
		tournament := matches.Eq(i).Find(".match-filler div div a").Eq(0).Text()
		uid := timestampStr + teamleft_text + teamright_text + tournament + "@lcalendar"
		flag := false // Ignore identical UIDs for now
		for j := 0; j < len(UIDs); j++ {
			if UIDs[j] == uid {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		UIDs = append(UIDs, uid)
		// Add event
		event := cal.AddEvent(uid)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(time.Unix(timestamp, 0))
		event.SetEndAt(time.Unix(timestamp+3600, 0))
		event.SetSummary(teamleft_text + " - " + teamright_text)
		event.SetLocation(tournament + " (" + matchFormat + ")")
	}
	return cal, nil
}
