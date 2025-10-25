package calendar

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/Napolitain/liquipedia_calendar/internal/handler"
	"github.com/PuerkitoBio/goquery"
)

func TestCreateCalendar(t *testing.T) {
	// Read all data from file resources/scrapping_test_data.html
	data, err := os.ReadFile("../../resources/scrapping_test_data_html")
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
	myQuery := handler.Query{Game: "starcraft2", Players: []string{"Maru", "Serral"}}
	cal, err := CreateCalendar(document, myQuery)
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
		iteratedEvents = append(iteratedEvents, events[i].Id())
	}
}
