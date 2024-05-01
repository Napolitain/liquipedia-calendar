package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"log"
	"os"
	"testing"
)

type QueriesMock struct {
	data []Query
}

func (queries QueriesMock) createCalendar(document *goquery.Document, player Query) (*ics.Calendar, error) {
	//TODO implement me
	panic("implement me")
}

func TestCreateCalendar(t *testing.T) {
	// Read all data from file resources/scrapping_test_data.html
	data, err := os.ReadFile("resources/scrapping_test_data.html")
	if err != nil {
		log.Fatal(err)
	}
	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	// Create iCalendar
	myQuery := Queries{data: []Query{{game: "starcraft2", players: []string{"Maru", "Serral"}}}}
	cal, err := myQuery.createCalendar(document, myQuery.data[0])
	if err != nil {
		t.Fatal(err)
	}
	// Check identical IDs
	events := cal.Events()
	var iteratedEvents []string
	for i := 0; i < len(events); i++ {
		for j := 0; j < len(iteratedEvents); j++ {
			if iteratedEvents[j] == events[i].Id() {
				t.Fatal("Identical ID between two events.")
			}
		}
		_ = append(iteratedEvents, events[i].Id())
	}
}
