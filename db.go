package main

import (
  "strings"

  "gopkg.in/mgo.v2"
)

// The DBConnection struct fields contain a session, collection, and boolean
// if the connection to the database is successful. Otherwise, the session
// and collection will be nil, and the boolean will be false.
type DBConnection struct {
  Session     *mgo.Session
  Collection  *mgo.Collection
  IsConnected bool
}

var db DBConnection

// Connects to a database given its connection string as long as there isn't
// already a connection established
func connect(connstr string) {
  if !db.IsConnected {
    if session, err := mgo.Dial(connstr); err != nil {
      db.IsConnected = false
      check("IGNORE", err)
    } else {
      session.SetMode(mgo.Monotonic, true)

      // Assign the fields to be used in the watcher
      db.Session = session
      db.Collection = session.DB("cs396").C("diffs")
      db.IsConnected = true

      // Ensure an index on the diff's timestamp and diff in the collection
      // where there are no duplicates and any duplicates already existing
      // are dropped
      index := mgo.Index {
        Key: []string{"sid", "timestamp"},
        Unique: true,
        DropDups: true,
      }
      err := db.Collection.EnsureIndex(index)
      check("PRINT", err)
    }
  }
}

// Inserts the given diff into the database
func insertDiff(diff Diff) {
  // Ignore logging duplicate key errors when saving to the
  // database
  err := db.Collection.Insert(&diff)
  if err != nil {
    if strings.Index(err.Error(), "duplicate key") < 0 {
      check("PRINT", err)
    }
  }
}