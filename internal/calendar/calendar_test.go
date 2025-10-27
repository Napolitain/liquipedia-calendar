package calendar

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
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

// TestUIDFormat verifies that UIDs are RFC 5545 compliant
func TestUIDFormat(t *testing.T) {
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

	events := cal.Events()
	if len(events) == 0 {
		t.Fatal("No events found in calendar")
	}

	// Log a sample UID to verify the format
	if len(events) > 0 {
		t.Logf("Sample UID: %q", events[0].Id())
	}

	// Verify each UID follows RFC 5545 compliant format
	for i, event := range events {
		uid := event.Id()

		// Check that UID is not empty
		if uid == "" {
			t.Errorf("Event %d has empty UID", i)
		}

		// Check that UID contains @lcalendar domain
		if !strings.Contains(uid, "@lcalendar") {
			t.Errorf("Event %d UID %q does not contain @lcalendar domain", i, uid)
		}

		// Check that UID does not contain spaces
		if strings.Contains(uid, " ") {
			t.Errorf("Event %d UID %q contains spaces", i, uid)
		}

		// Check that UID follows format: unique-id@domain
		// unique-id should only contain alphanumeric and limited special characters
		parts := strings.Split(uid, "@")
		if len(parts) != 2 {
			t.Errorf("Event %d UID %q does not follow unique-id@domain format", i, uid)
			continue
		}

		uniqueID := parts[0]
		domain := parts[1]

		// Verify unique ID contains only hexadecimal characters (since we use SHA256)
		hexPattern := regexp.MustCompile("^[0-9a-f]+$")
		if !hexPattern.MatchString(uniqueID) {
			t.Errorf("Event %d UID unique part %q is not a valid hex string", i, uniqueID)
		}

		// Verify domain is correct
		if domain != "lcalendar" {
			t.Errorf("Event %d UID domain %q is not 'lcalendar'", i, domain)
		}

		// Verify the unique ID has reasonable length (SHA256 produces 64 hex chars)
		if len(uniqueID) != 64 {
			t.Errorf("Event %d UID unique part has length %d, expected 64 (SHA256 hash)", i, len(uniqueID))
		}
	}
}
