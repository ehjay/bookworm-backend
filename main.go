package main

import (
    "encoding/json"
    "log"
    "fmt"
    "net/http"
    "time"

    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/gorilla/context"
)

type book struct {
  ID      bson.ObjectId `json:"id" bson:"_id"`
  Title   string        `json: "title" bson: "title"`
  Author  string        `json: "author" bson: "author"`
  Created time.Time     `json:"created" bson:"created"`
}

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
  var books []*book
  if err := db.DB("bookworm").C("books").
    Find(nil).Sort("-when").Limit(100).All(&books); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

  if err := json.NewEncoder(w).Encode(books); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
}

func addBook(w http.ResponseWriter, r *http.Request) {
    db := context.Get(r, "database").(*mgo.Session)
    var b book
    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }
    b.ID = bson.NewObjectId()
    b.Created = time.Now()
    if err := db.DB("bookworm").C("books").Insert(&b); err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }
}

func main() {
    db, err := mgo.Dial("localhost")
    if err != nil {
      log.Fatal("cannot dial mongo", err)
    }
    defer db.Close() // clean up when we’re done
    mux := goji.NewMux()
    mux.Use(logging)
    mux.Use(withDB(db))
    mux.HandleFunc(pat.Get("/books"), allBooks)
    mux.HandleFunc(pat.Post("/book"), addBook)
    fmt.Println("Server starting...");
    if err := http.ListenAndServe("localhost:8080", mux); err != nil {
      log.Fatal(err)
    }
}
