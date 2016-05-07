package main

import (
  "log"
  "math"
  "path/filepath"
  "os"
  "strings"
  "time"

  "github.com/fsnotify/fsnotify"
)

// Error checking handler where action indicates how to handle the error.
// More cases can be added to handle errors in different ways.
func check(action string, err error) {
  if err != nil {
    switch action {
      // Ignore the error
      case "IGNORE":
      // Log the error
      case "PRINT":
        log.Println("ERROR:", err)
    }
  }
}

// The current time formatted as RFC1123Z
func currentTime() string {
  return time.Now().Format(time.RFC1123Z)
}

// Gets the size of a file based on its properties given its path
func fileSize(path string) float64 {
  if info, err := os.Stat(path); err != nil {
    return float64(0)
  } else {
    return float64(info.Size())
  }
}

// Source: http://goo.gl/TFeP7R
// Wraps IsDir() to take account of the existence of the given path
// rather than simply returning an error
func isDirectory(path string) bool {
  if file, err := os.Stat(path); os.IsNotExist(err) {
    // Path does not exist, thus it's not a directory
    return false
  } else {
    return file.IsDir()
  }
}

// Gets the size difference (in bytes) between the old and new paths
func sizeDiff(oldPath string, newPath string) int {
  oldSize, newSize := fileSize(oldPath), fileSize(newPath)

  return int(math.Abs(newSize - oldSize))
}

// Source: http://goo.gl/O6kkq2
// Walks the directory from a given path and watches each subdirectory
// while ignoring the .tmp directory and its subdirectories
func watchDirectory(watcher *fsnotify.Watcher, path string) {
  filepath.Walk(path, func(p string, _ os.FileInfo, _ error) error {
    if isDirectory(p) && strings.Index(p, ".tmp") < 0 {
      if err := watcher.Add(p); err != nil {
        check("PRINT", err)
      }
    }

    return nil
  })
}

// Source: http://goo.gl/XfolFX
// Writes to or creates a file given its path with a given text
func writeOrCreateFile(path string, text string) {
  // Open file as write-only, or create if it doesn't exist
  file, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0644)
  check("PRINT", err)

  defer file.Close()

  _, err = file.WriteString(text)
  check("PRINT", err)
}