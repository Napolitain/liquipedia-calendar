package main

import (
	"fmt"
	"google.golang.org/appengine"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		response, err := getFromLiquipedia("ageofempires")
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}
		if response.StatusCode != 200 {
			_, _ = fmt.Fprint(w, response.Status)
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		err = saveToCache(r.Context(), string(body[:]))
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
		}
		_, _ = fmt.Fprint(w, string(body[:])+" generated.")
	} else {
		_, _ = fmt.Fprint(w, string(item.Value[:])+" from memcached.")
	}
}
