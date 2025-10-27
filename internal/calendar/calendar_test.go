package calendar

import (
	"bytes"
	"fmt"
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
	// Read all data from file resources/scraping_test_data.html
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

// TestRFC5545Compliance verifies full RFC 5545 compliance for calendar and events
func TestRFC5545Compliance(t *testing.T) {
	// Read all data from file resources/scraping_test_data.html
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

	// Test Calendar-level RFC 5545 required properties
	t.Run("Calendar Required Properties", func(t *testing.T) {
		calData := cal.Serialize()

		// VERSION:2.0 is required
		if !strings.Contains(calData, "VERSION:2.0") {
			t.Error("Calendar missing required VERSION:2.0 property")
		}

		// PRODID is required
		if !strings.Contains(calData, "PRODID:") {
			t.Error("Calendar missing required PRODID property")
		}

		// CALSCALE (optional but recommended)
		if !strings.Contains(calData, "CALSCALE:GREGORIAN") {
			t.Error("Calendar missing recommended CALSCALE:GREGORIAN property")
		}

		// METHOD (optional)
		if !strings.Contains(calData, "METHOD:") {
			t.Error("Calendar missing METHOD property")
		}
	})

	// Test Event-level RFC 5545 compliance
	events := cal.Events()
	if len(events) == 0 {
		t.Fatal("No events found in calendar for compliance testing")
	}

	for i, event := range events {
		t.Run(fmt.Sprintf("Event %d RFC 5545 Compliance", i), func(t *testing.T) {
			// Required properties per RFC 5545 Section 3.6.1

			// UID is required
			uid := event.Id()
			if uid == "" {
				t.Error("Event missing required UID property")
			}

			// DTSTAMP is required
			_, err := event.GetDtStampTime()
			if err != nil {
				t.Errorf("Event missing required DTSTAMP property: %v", err)
			}

			// DTSTART is required
			_, err = event.GetStartAt()
			if err != nil {
				t.Errorf("Event missing required DTSTART property: %v", err)
			}

			// DTEND or DURATION is required (we use DTEND)
			_, err = event.GetEndAt()
			if err != nil {
				t.Errorf("Event missing required DTEND property: %v", err)
			}

			// Optional but recommended properties
			summaryProp := event.GetProperty("SUMMARY")
			if summaryProp == nil || summaryProp.Value == "" {
				t.Error("Event missing recommended SUMMARY property")
			}

			locationProp := event.GetProperty("LOCATION")
			if locationProp == nil || locationProp.Value == "" {
				t.Error("Event missing recommended LOCATION property")
			}

			descriptionProp := event.GetProperty("DESCRIPTION")
			if descriptionProp == nil || descriptionProp.Value == "" {
				t.Error("Event missing recommended DESCRIPTION property")
			}

			statusProp := event.GetProperty("STATUS")
			if statusProp == nil || statusProp.Value == "" {
				t.Error("Event missing recommended STATUS property")
			}

			sequenceProp := event.GetProperty("SEQUENCE")
			if sequenceProp == nil {
				t.Error("Event missing recommended SEQUENCE property")
			}
		})
	}
}

// TestCalendarSerializationFormat verifies the calendar serializes correctly
func TestCalendarSerializationFormat(t *testing.T) {
	// Read all data from file resources/scraping_test_data.html
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

	// Serialize calendar
	calData := cal.Serialize()

	// Check basic iCalendar structure
	if !strings.HasPrefix(calData, "BEGIN:VCALENDAR") {
		t.Error("Calendar serialization should start with BEGIN:VCALENDAR")
	}

	if !strings.HasSuffix(strings.TrimSpace(calData), "END:VCALENDAR") {
		t.Error("Calendar serialization should end with END:VCALENDAR")
	}

	// Check for VEVENT blocks
	if !strings.Contains(calData, "BEGIN:VEVENT") {
		t.Error("Calendar should contain at least one BEGIN:VEVENT")
	}

	if !strings.Contains(calData, "END:VEVENT") {
		t.Error("Calendar should contain at least one END:VEVENT")
	}

	// Log sample for verification
	lines := strings.Split(calData, "\n")
	t.Logf("Calendar has %d lines", len(lines))
	if len(lines) > 0 {
		t.Logf("First line: %s", lines[0])
		t.Logf("Last line: %s", strings.TrimSpace(lines[len(lines)-1]))
	}
}
