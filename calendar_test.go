package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"google.golang.org/appengine/aetest"
	"log"
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
	return // app engine context is difficult to setup
	// Create context
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// Get data
	data, err := getData(ctx, "starcraft2")
	if err != nil {
		t.Fatal(err)
	}
	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	// Create iCalendar
	myQuery := Queries{}
	cal, err := myQuery.createCalendar(document, "")
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
