package main

import (
  "encoding/json"
)

// Stores data for each diff generated
type Diff struct {
  StudentID int    `json:"sid`
  Timestamp string `json:"timestamp"`
  Op        string `json:"op"`
  Path      string `json:"path"`
  SizeDiff  int    `json:"sizediff"`
  Diff      string `json:"diff"`
}

// Stores every diff to be sent off to the database
type Diffs struct {
  Diffs []Diff `json:"diffs"`
}

// Keep this global so it doesn't get reinitialized
var diffs Diffs

// Takes diffs and appends every diff passed to this function.
// diffs is then JSON encoded with 2-space indentation for readability and
// written to the given path.
func writeJSONToPath(diff Diff, path string) {
  diffs.Diffs = append(diffs.Diffs, diff)
  // JSON encode with 2-space indentation for readability
  jDiffs, err := json.MarshalIndent(diffs, "", "  ")
  check("PRINT", err)

  writeOrCreateFile(path, string(jDiffs))
}