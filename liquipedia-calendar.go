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
	"strconv"
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

	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create iCalendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//liquipedia-calendar//en")
	cal.SetVersion("2.0")

	matches := document.Find(".infobox_matches_content")
	for i := 0; i < matches.Size(); i++ {
		teamleft := matches.Eq(i).Find(".team-left a").Eq(0).Text()
		teamright := matches.Eq(i).Find(".team-right span:not(.flag):not(.team-template-image):not(.team-template-team-short) a").Eq(0).Text()
		matchFormat := matches.Eq(i).Find(".versus abbr").Eq(0).Text()
		timestampStr, exist := matches.Eq(i).Find("[data-timestamp]").Eq(0).Attr("data-timestamp")
		if exist != true {
			log.Fatal("Timestamp doesn't exist.")
			return
		}
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			log.Fatal("Error while converting string to int.")
			return
		}
		tournament := matches.Eq(i).Find(".match-filler div div a").Eq(0).Text()
		uid := timestampStr + teamleft + teamright + tournament + "@lcalendar"

		event := cal.AddEvent(uid)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(time.Unix(timestamp, 0))
		event.SetEndAt(time.Unix(timestamp+3600, 0))
		event.SetSummary(teamleft + " - " + teamright)
		event.SetLocation(tournament + " (" + matchFormat + ")")
		if i == 0 {
			log.Println(event.Serialize())
		}
	}
	w.Header().Set("Content-Disposition", "attachment; filename=sc2calendar.ics")
	w.Header().Set("Content-Type", "text/calendar")
	_, err = fmt.Fprintf(w, cal.Serialize())
	if err != nil {
		log.Fatal("Error while printing serialized calendar.")
		return
	}
}
