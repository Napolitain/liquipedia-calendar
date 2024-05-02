package main

import (
	"bytes"
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "", 0)

func main() {
	// Creates a client.
	ctx := context.Background()
	client, err := logging.NewClient(ctx, "liquipedia-calendar")
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}
	defer client.Close()
	logger = client.Logger("main-service").StandardLogger(logging.Info)

	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Printf("Defaulting to port %s", port)
	}

	logger.Printf("Listening on port %s", port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is for the root path. If not, return 404.
	if r.URL.Path != "/" {
		logger.Println("Path not supported." + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	// Get query string's name from querystring.
	querystring := r.URL.Query().Get("query")
	if querystring == "" {
		logger.Println("No query string provided.")
		return
	}

	// Get query struct
	queries := newQueries(querystring)

	// Get from cache the game+player calendar if cached. (Superstar player case).
	calendar, err := getFromCachePlayer(r.Context(), queries.data[0].game, queries.data[0].players[0])
	if err == nil {
		sendCalendar(w, err, string(calendar.Value))
		return
	}

	// Get data from either cache (game generic case) or scrapping. JSON already parsed and filtered HTML.
	data, fromCache, err := getData(r.Context(), queries.data[0].game) // TODO: Handle multiple games
	if err != nil {
		logger.Println(err)
		return
	}

	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		logger.Println(err)
	}

	// Create iCalendar
	cal, err := queries.createCalendar(document, queries.data[0])
	if err != nil {
		logger.Println(err)
	}

	serializedCalendar := cal.Serialize()

	// If it is for a single player, save to cache the game+player calendar (superstar player case).
	err = saveToCachePlayer(r.Context(), serializedCalendar, queries.data[0].game, queries.data[0].players[0])
	sendCalendar(w, err, serializedCalendar)
}

// sendCalendar function is used to send the calendar to the user.
// It sends the calendar in iCalendar format.
func sendCalendar(w http.ResponseWriter, err error, serializedCalendar string) {
	if err != nil {
		logger.Println(err)
	}

	// Set headers for the response
	w.Header().Set("Content-Disposition", "attachment; filename=liquipedia.ics")
	w.Header().Set("Content-Type", "text/calendar")

	_, err = fmt.Fprintf(w, serializedCalendar)
	if err != nil {
		logger.Println("Error while printing serialized calendar.")
		return
	}
}
