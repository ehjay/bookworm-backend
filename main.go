package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "goji.io"
    "goji.io/pat"
)

func allBooks(w http.ResponseWriter, r *http.Request) {
    jsonOut, _ := json.Marshal("the books")
    fmt.Fprintf(w, string(jsonOut))
}

func logging(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Received request: %v\n", r.URL)
        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

func main() {
    mux := goji.NewMux()
    mux.HandleFunc(pat.Get("/books"), allBooks)
    mux.Use(logging)
    fmt.Println("Starting server...");
    http.ListenAndServe("localhost:8080", mux)
}
