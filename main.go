package main

import (
  "io/ioutil"
  "log"
  "os"

  "github.com/BurntSushi/toml"
)

// Contains configuration settings set by the user in config.toml
type Config struct {
  Class       string
  Assignment  string
  Directories []string
  ConnStr     string
}

var conf Config
var snapshot string

// Deletes the .tmp directory to begin a new session. Read the config.toml
// configuration file to retrieve all user-specified settings. A snapshot of
// the current directory is then created, then the watcher is initialized.
func main() {
  if isDirectory(".tmp") {
    os.RemoveAll(".tmp")
  }

  log.Println("Reading configurations...")

  // Parse config.toml to retrieve configurations for the watcher.
  // If config.toml doesn't exist, continue on with execution as it's not
  // absolutely necessary.
  tData, err := ioutil.ReadFile("config.toml")
  check("PRINT", err)

  _, err = toml.Decode(string(tData), &conf)
  check("PRINT", err)

  log.Println("Configurations set:")
  log.Println("  CLASS:       ", conf.Class)
  log.Println("  ASSIGNMENT:  ", conf.Assignment)
  log.Println("  DIRS WATCHED:", conf.Directories)

  // Generate a snapshot of the original directory
  snapshot = currentTime()
  createSnapshot("./", ".tmp/snapshots/" + snapshot)

  bigBrotherIsWatching()
}