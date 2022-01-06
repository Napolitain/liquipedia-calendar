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
	if err != nil {
		s := save(r.Context())
		_, _ = fmt.Fprint(w, s+" generated.")
	} else {
		_, _ = fmt.Fprint(w, string(item.Value[:])+" from memcached.")
	}
}

// save function is used to save data on Memcached server.
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

// get function is used to retrieve data from Memcached server.
func get(ctx context.Context) (*memcache.Item, error) {
	return memcache.Get(ctx, "data")
}
