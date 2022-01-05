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

	item, err := get(r.Context())
	if err != nil && item != nil {
		fmt.Fprint(w, string(item.Value[:])+" from memcached.")
	} else {
		s := save(r.Context())
		fmt.Fprint(w, s+" generated.")
	}
}

func save(ctx context.Context) string {
	now := time.Now()
	item := &memcache.Item{
		Key:   "data",
		Value: []byte(now.String()),
	}
	err := memcache.Set(ctx, item)
	if err != nil {
		return err.Error()
	}
	return now.String()
}

func get(ctx context.Context) (*memcache.Item, error) {
	return memcache.Get(ctx, "data")
}
