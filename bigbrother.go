// Source: https://godoc.org/github.com/fsnotify/fsnotify

package main

import (
  "log"
  "strings"
  "time"

  "github.com/fsnotify/fsnotify"
)

/* bigBrotherIsWatching() initializes a new watcher to start watching
 * directories. It attempts a first connection to the database. A goroutine then
 * starts to reattempt connecting periodically regardless of whether the it has
 * been established or not.
 * The same goroutine has the watcher listening for events and errors. When an
 * event is received, it begins to generate the diff associate with the path(s)
 * of the event. The diff is stored in .tmp/output.json regardless if there is
 * a connection to the database. When a connection is established, each diff in
 * .tmp/output.json will be inserted into the database. The database is set up
 * to prevent duplicate documents. If there is already a connection, the diff
 * is inserted immediately.
 * The event listening ends with creating a snapshot of the directory at its
 * current state and storing it in .tmp/snapshots.
 * At the end of the goroutine and watcher, they go through all user-specified
 * directories to traverse their file tree and add watcher for any directory
 * discovered. If there are no user-specified directories, the entire project
 * directory will be watched (including all subdirectories).
 */
func bigBrotherIsWatching() {
  watcher, err := fsnotify.NewWatcher()
  check("PRINT", err)

  defer watcher.Close()

  // First attempt at connecting to the database
  if connect(conf.ConnStr); db.IsConnected {
    log.Println("Connected to database")

    defer db.Session.Close()
  } else {
    log.Println("Unable to reach the database")
  }

  // diffID is a temporary solution for uniquely identifying diffs.
  // This should be replaced with some sort of identifier for each student.
  diffID := 0

  ssDir := ".tmp/snapshots/"
  // Attempt to reconnect to the database every 10 seconds
  ticker := time.NewTicker(10 * time.Second)
  done := make(chan bool)
  go func() {
    for {
      select {
        case <- ticker.C:
          // Attempt reconnection
          if connect(conf.ConnStr); db.IsConnected {
            defer db.Session.Close()

            // Insert every saved diff into the database
            for _, diff := range diffs.Diffs {
              insertDiff(diff)
            }
          }

        case event := <- watcher.Events:
          // Parse the event name to get the operation that occurred
          index := strings.Index(event.String(), ":") + 2
          op := event.String()[index:]
          switch op {
            // Ignore permission changes
            case "CHMOD":
            default:
               // Path of file/directory
              evDir := event.Name
              // Stop watching a directory when it's deleted
              if op == "REMOVE" {
                watcher.Remove(evDir)
              }

              // Generate the diff with the number of bytes changed and the
              // changes in the old and new files
              size := sizeDiff(ssDir + snapshot + "/" + evDir, evDir)
              diff := diff(ssDir + snapshot + "/" + evDir, evDir)
              jDiff := Diff {
                StudentID: diffID,
                Timestamp: time.Now().String(),
                Op:        op,
                Path:      evDir,
                SizeDiff:  size,
                Diff:      diff,
              }

              diffID += 1

              // Output the diff to .tmp/output.json and insert it into the
              // database if there is a connection
              if db.IsConnected {
                writeJSONToPath(jDiff, ".tmp/output.json")
                insertDiff(jDiff)
              } else {
                writeJSONToPath(jDiff, ".tmp/output.json")
              }

              snapshot = currentTime()
              createSnapshot("./", ssDir + snapshot)

              log.Printf("EVENT: %s \n", event)
          }

          // Go through each specified directory and watch all new directories,
          // or the entire directory if none were specified
          if len(conf.Directories) > 0 {
            for _, dir := range conf.Directories {
              watchDirectory(watcher, dir)
            }
          } else {
            watchDirectory(watcher, "./")
          }

        case err := <- watcher.Errors:
          log.Printf("ERROR: %s \n", err)
      }
    }
  }()

  // Watch each directory specified in config.toml. Skip the directories
  // that don't exist, and watch the entire directory if none were specified,
  // or if all specified directories don't exist (this shouldn't happen)
  if len(conf.Directories) > 0 {
    errCounter := 0

    // Count the number of directories that don't exist
    for _, dir := range conf.Directories {
      if err = watcher.Add(dir); err != nil {
        errCounter += 1
      }
    }

    // Watch the entire directory since all specified ones didn't exist
    if errCounter == len(conf.Directories) {
      err = watcher.Add("./")
      check("PRINT", err)
    }
  } else {
    err = watcher.Add("./")
    check("PRINT", err)
  }

  log.Println("Big Brother is watching you...")
  <- done
}