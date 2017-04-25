package main

import (
    "encoding/json"
    "log"
    "fmt"
    "net/http"

    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"

    "github.com/gorilla/context"
)

// middleware

func logging(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Received request: %v\n", r.URL)
        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

type Adapter func(http.Handler) http.Handler

func withDB(db *mgo.Session) Adapter {
  return func(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
      dbsession := db.Copy()
      defer dbsession.Close() // clean up
      context.Set(r, "database", dbsession)
      h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
  }
}

// routes

func allBooks(w http.ResponseWriter, r *http.Request) {
    db := context.Get(r, "database").(*mgo.Session)
    fmt.Println(db);
    jsonOut, _ := json.Marshal("the books")
    fmt.Fprintf(w, string(jsonOut))
}

func main() {
    db, err := mgo.Dial("localhost")
    if err != nil {
      log.Fatal("cannot dial mongo", err)
    }
    defer db.Close() // clean up when we’re done

    mux := goji.NewMux()

    mux.Use(logging)
    // mux.Use(withDB(db))

    mux.HandleFunc(pat.Get("/books"), allBooks)

    fmt.Println("Server starting...");
    if err := http.ListenAndServe("localhost:8080", mux); err != nil {
      log.Fatal(err)
    }
}
