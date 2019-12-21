package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/cors"
)

// SearchRequest structure
type SearchRequest struct {
	Prefixed  string `json:"prefix"`
	Confirmed bool   `json:"confirmed"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello!\n")
}

func search(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t SearchRequest
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	result, err := FindWord(t)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	wordsFound, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(wordsFound)
}

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
	}
	mux := http.NewServeMux()

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/search", search)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
