package main

import (
	"bytes"
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"k8s.io/apimachinery/pkg/util/sets"
	_ "k8s.io/apimachinery/pkg/util/sets"
	"log"
	"net/http"
	"os"
	"strings"
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

	// Get from cache the queries calendar if cached. (Superstar player case).
	calendar, err := getFromCachePlayer(r.Context(), r.URL.Query().Get("query"))
	if err == nil {
		sendCalendar(w, err, string(calendar.Value))
		return
	}

	// Get query string's name from querystring.
	querystring := r.URL.Query().Get("query")
	if querystring == "" {
		logger.Println("No query string provided.")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get query struct
	queries, err := newQueries(querystring)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// If the game inside query is not valid, return bad request.
	if !isValidGame(r.Context(), queries.data[0].game) {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get data from either cache (game generic case) or scrapping. JSON already parsed and filtered HTML.
	data, err := getData(r.Context(), queries.data[0].game) // TODO: Handle multiple games
	if err != nil {
		logger.Println(err)
		http.NotFound(w, r)
		return
	}

	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// Create iCalendar
	cal, err := queries.createCalendar(document, queries.data[0])
	if err != nil {
		logger.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	serializedCalendar := cal.Serialize()

	// If it is for a single player, save to cache the game+player calendar (superstar player case).
	err = saveToCachePlayer(r.Context(), serializedCalendar, r.URL.Query().Get("query"))
	sendCalendar(w, err, serializedCalendar)
}

// TODO: do the function, test and move to other file.
func isValidGame(ctx context.Context, game string) bool {
	// List of games supported by Liquipedia API : important to avoid not only errors, but attacks.
	// Retrieve from the cache first.
	item, err := getGamesFromCache(ctx)
	// Cache result
	if err == nil {
		// Convert string to sets.String
		gamesMap := sets.NewString()
		// Deserialize string1,string2 to []string then to Map
		gamesMap.Insert(strings.Split(string(item.Value), ",")...)
		// Check if the game is in the map, return true if it is.
		return gamesMap.Has(game)
	}
	// Cache miss, log the issue, but continue.
	logger.Println(err)
	// Otherwise, from the DB.
	// Otherwise, from the API.
	gamesMap, err := fetchGames()
	if err != nil {
		logger.Println(err)
		return false
	}
	// Save to cache
	// Serialize Map to []string, then to string1,string2
	gamesList := gamesMap.List()
	games := strings.Join(gamesList, ",")
	err = saveGamesToCache(ctx, games)
	return gamesMap.Has(game)
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
