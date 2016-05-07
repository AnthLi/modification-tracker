// Source: https://godoc.org/github.com/pmezard/go-difflib/difflib

package main

import(
  "io/ioutil"
  "math"
  "strings"

  "github.com/pmezard/go-difflib/difflib"
)

// Get the differences between two files and return the output.
// If the oldPath is an empty string, then newPath indicates a new file.
// Otherwise, oldPath is the original file and newPath is the modified file.
func diff(oldPath string, newPath string) string {
  // Read both files as strings to figure out which file has more lines
  readOld, _ := ioutil.ReadFile(oldPath)
  readNew, _ := ioutil.ReadFile(newPath)
  oldLines :=   string(readOld[:])
  newLines :=   string(readNew[:])
  oldLength :=  float64(len(strings.Split(oldLines, "\n")))
  newLength :=  float64(len(strings.Split(newLines, "\n")))

  // Create the diff
  var diff difflib.UnifiedDiff
  if len(oldPath) < 1 {
    // The file is new
    diff = difflib.UnifiedDiff {
      B:       difflib.SplitLines(newLines),
      ToFile:  newPath,
      Context: int(newLength),
      Eol:     "\n",
    }
  } else {
    // The file is modified
    diff = difflib.UnifiedDiff {
      A:        difflib.SplitLines(oldLines),
      B:        difflib.SplitLines(newLines),
      FromFile: oldPath,
      ToFile:   newPath,
      Context:  int(math.Max(oldLength, newLength)),
      Eol:      "\n",
    }
  }

  diffString, _ := difflib.GetUnifiedDiffString(diff)
  return diffString
}