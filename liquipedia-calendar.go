package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get game's name from querystring.
	game := r.URL.Query().Get("game")
	if game == "" {
		log.Fatal("No query string provided.")
		return
	}

	// Get data from either cache or scrapping. JSON already parsed and filtered HTML.
	data, err := getData(r.Context(), game)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create iCalendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//liquipedia-calendar//en")
	cal.SetVersion("3.0")
	event := cal.AddEvent(fmt.Sprintf("id@domain", p.SessionKey.IntID()))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(time.Now())
	event.SetEndAt(time.Now())
	event.SetSummary("Summary")
	event.SetLocation("Address")
	event.SetDescription("Description")
	event.SetURL("https://URL/")
	event.AddRrule(fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", time.Now().Month(), time.Now().Day()))
	event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
	event.AddAttendee("reciever or participant", ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))
	// cal.Serialize()

	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	document.Find(".left-content article .post-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
