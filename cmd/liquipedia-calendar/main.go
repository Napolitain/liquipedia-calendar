package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Napolitain/liquipedia_calendar/internal/cache"
	icalendar "github.com/Napolitain/liquipedia_calendar/internal/calendar"
	"github.com/Napolitain/liquipedia_calendar/internal/handler"
	"github.com/Napolitain/liquipedia_calendar/internal/scraping"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"k8s.io/apimachinery/pkg/util/sets"
)

func main() {
	// Create structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize cache
	cache.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Info("Defaulting to port", "port", port)
	}

	slog.Info("Starting server", "port", port)
	
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		indexHandler(ctx)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), requestHandler); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	// Check if the request is for the root path. If not, return 404.
	if string(ctx.Path()) != "/" {
		slog.Info("Path not supported", "path", string(ctx.Path()))
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}

	// Get from cache the queries calendar if cached. (Superstar player case).
	queryString := string(ctx.QueryArgs().Peek("query"))
	calendar, err := cache.GetPlayerCalendar(queryString)
	if err == nil && calendar != "" {
		sendCalendar(ctx, calendar)
		return
	}

	// Get query string's name from querystring.
	if queryString == "" {
		slog.Info("No query string provided")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	// Get query struct
	queries, err := handler.NewQueries(queryString)
	if err != nil {
		slog.Error("Failed to parse queries", "error", err)
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	// If the game inside query is not valid, return bad request.
	if !isValidGame(queries.Data[0].Game) {
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	// Get data from either cache (game generic case) or scrapping. JSON already parsed and filtered HTML.
	data, err := getData(queries.Data[0].Game)
	if err != nil {
		slog.Error("Failed to get data", "error", err)
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}

	// Load the HTML document
	reader := bytes.NewReader(data)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		slog.Error("Failed to parse HTML", "error", err)
		ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
		return
	}

	// Create iCalendar
	cal, err := icalendar.CreateCalendar(document, queries.Data[0])
	if err != nil {
		slog.Error("Failed to create calendar", "error", err)
		ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
		return
	}

	serializedCalendar := cal.Serialize()

	// If it is for a single player, save to cache the game+player calendar (superstar player case).
	cache.SetPlayerCalendar(queryString, serializedCalendar)
	sendCalendar(ctx, serializedCalendar)
}

func isValidGame(game string) bool {
	// List of games supported by Liquipedia API : important to avoid not only errors, but attacks.
	// Retrieve from the cache first.
	cachedGames := cache.GetGames()
	// Cache hit
	if cachedGames != "" {
		// Convert string to sets.String
		gamesMap := sets.NewString()
		// Deserialize string1,string2 to []string then to Map
		gamesMap.Insert(strings.Split(cachedGames, ",")...)
		// Check if the game is in the map, return true if it is.
		return gamesMap.Has(game)
	}
	// Cache miss, log the issue, but continue.
	slog.Info("Games cache miss, fetching from API")
	// Otherwise, from the API.
	gamesMap, err := scraping.FetchGames()
	if err != nil {
		slog.Error("Failed to fetch games", "error", err)
		return false
	}
	// Save to cache
	// Serialize Map to []string, then to string1,string2
	gamesList := gamesMap.List()
	games := strings.Join(gamesList, ",")
	cache.SetGames(games)
	return gamesMap.Has(game)
}

// sendCalendar function is used to send the calendar to the user.
// It sends the calendar in iCalendar format.
func sendCalendar(ctx *fasthttp.RequestCtx, serializedCalendar string) {
	// Set headers for the response
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=liquipedia.ics")
	ctx.Response.Header.Set("Content-Type", "text/calendar")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(serializedCalendar)
}

// getData returns data in []byte format from either cache or scrapping
func getData(game string) ([]byte, error) {
	// Get data from cache server
	cachedData := cache.GetGameData(game)
	if cachedData != nil {
		return cachedData, nil
	}

	// If fail, get data from scrapping
	data, err := scraping.GetFromLiquipedia(game)
	if err != nil {
		slog.Error("Failed to get data from Liquipedia", "error", err)
		return nil, err
	}

	// Save to cache server
	cache.SetGameData(game, data)
	return data, nil
}
