// This needs work/replacement since JSON was chosen over CSV early on during
// implementation.

package main

import(
  "encoding/csv"
  "os"
)

// Output the given CSV to the given path
func writeCSVToFile(csvData [][]string, path string) {
  outputFile, err := os.Create(path)
  check("PRINT", err)

  defer outputFile.Close()

  writer := csv.NewWriter(outputFile)

  for _, data := range csvData {
    err = writer.Write(data)
    check("PRINT", err)
  }

  writer.Flush()
}