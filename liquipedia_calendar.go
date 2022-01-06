package main

import (
	"context"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
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

	item, err := getFromCache(r.Context())
	if err != nil {
		cache, err := saveToCache(r.Context())
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
		}
		_, _ = fmt.Fprint(w, cache.String()+" generated.")
	} else {
		_, _ = fmt.Fprint(w, string(item.Value[:])+" from memcached.")
	}
}

// saveToCache function is used to saveToCache data on Memcached server.
func saveToCache(ctx context.Context) (time.Time, error) {
	now := time.Now()
	item := &memcache.Item{
		Key:   "data",
		Value: []byte(now.String()),
	}
	err := memcache.Set(ctx, item)
	return now, err
}

// getFromCache function is used to retrieve data from Memcached server.
func getFromCache(ctx context.Context) (*memcache.Item, error) {
	return memcache.Get(ctx, "data")
}

// getFromLiquipedia function
func getFromLiquipedia(game string) (*http.Response, error) {
	client := http.Client{}
	url := "https://liquipedia.net/" + game + "api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Host":            []string{"liquipedia.net"},
		"Content-Type":    []string{"application/json"},
		"Accept-Encoding": []string{"gzip"},
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
