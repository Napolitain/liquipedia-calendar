package main

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"io/ioutil"
	"strconv"
	"time"
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

func createCalendar(document *goquery.Document) (*ics.Calendar, error) {
	// Create iCalendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//liquipedia-calendar//en")
	cal.SetVersion("2.0")
	cal.SetLastModified(time.Now())
	// Create events
	matches := document.Find(".infobox_matches_content")
	for i := 0; i < matches.Size(); i++ {
		// Get event info
		teamleft := matches.Eq(i).Find(".team-left a").Eq(0).Text()
		teamright := matches.Eq(i).Find(".team-right span:not(.flag):not(.team-template-image):not(.team-template-team-short) a").Eq(0).Text()
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
		uid := timestampStr + teamleft + teamright + tournament + "@lcalendar"
		// Add event
		event := cal.AddEvent(uid)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(time.Unix(timestamp, 0))
		event.SetEndAt(time.Unix(timestamp+3600, 0))
		event.SetSummary(teamleft + " - " + teamright)
		event.SetLocation(tournament + " (" + matchFormat + ")")
	}
	return cal, nil
}
