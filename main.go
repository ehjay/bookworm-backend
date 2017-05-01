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

type Book struct {
  ID bson.ObjectId `json:"id" bson:"_id"`
  Title string `json:"title" bson:"title"`
  Author string `json:"author" bson:"author"`
  Created time.Time `json:"created" bson:"created"`
}

// middleware

func cors(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
        if (r.Method == "OPTIONS") {
          w.WriteHeader(200)
        } else {
          h.ServeHTTP(w, r)
        }
    }
    return http.HandlerFunc(fn)
}

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
  var books []*Book
  if err := db.DB("bookworm").C("books").
    Find(nil).Sort("title").Limit(100).All(&books); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

  jsonOut, _ := json.Marshal(books)
  fmt.Fprintf(w, string(jsonOut))
}

func addBook(w http.ResponseWriter, r *http.Request) {
    db := context.Get(r, "database").(*mgo.Session)
    var b Book
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

func removeBook(w http.ResponseWriter, r *http.Request) {
    db := context.Get(r, "database").(*mgo.Session)
    id := pat.Param(r, "id")
    log.Print("id: " + id)
    if err := db.DB("bookworm").C("books").Remove(bson.M{"_id": bson.ObjectIdHex(id)}); err != nil {
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
    mux.Use(cors)
    mux.Use(logging)
    mux.Use(withDB(db))
    mux.HandleFunc(pat.Get("/books"), allBooks)
    mux.HandleFunc(pat.Post("/book"), addBook)
    mux.HandleFunc(pat.Delete("/book/:id"), removeBook)
    fmt.Println("Server starting...");
    if err := http.ListenAndServe("localhost:8080", mux); err != nil {
      log.Fatal(err)
    }
}
