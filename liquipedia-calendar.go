package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandler)
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

	// Get query string's name from querystring.
	querystring := r.URL.Query().Get("query")
	if querystring == "" {
		log.Fatal("No query string provided.")
		return
	}

	// Get query struct
	queries := newQueries(querystring)

	// Get data from either cache or scrapping. JSON already parsed and filtered HTML.
	data, err := getData(r.Context(), querystring)
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
	cal, err := queries.createCalendar(document)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=liquipedia.ics")
	w.Header().Set("Content-Type", "text/calendar")
	_, err = fmt.Fprintf(w, cal.Serialize())
	if err != nil {
		log.Fatal("Error while printing serialized calendar.")
		return
	}
}
