// Source for copyFile and copyDir: https://goo.gl/Gl3NnQ

// The functions copyFile, copyDir, and createSnapshot are used to create
// the snapshot of a directory. The snapshots are used to track the user's
// progress with their work.

package main

import (
  "io"
  "os"
)

// Copies a source file to a destination file
func copyFile(src string, dst string) {
  srcFile, err := os.Open(src)
  check("PRINT", err)

  defer srcFile.Close()

  dstFile, err := os.Create(dst)
  check("PRINT", err)

  defer dstFile.Close()

  // Copy the source file to the destination file and retain persmissions
  // for the new file
  if _, err := io.Copy(dstFile, srcFile); err == nil {
    if srcInfo, err := os.Stat(src); err == nil {
      err = os.Chmod(dst, srcInfo.Mode())
      check("PRINT", err)
    }
  }
}

// Copies a source directory to a destination directory
func copyDir(src string, dst string) {
  srcInfo, err := os.Stat(src)
  check("PRINT", err)

  // Retain permissions of the source when creating the destination directory
  err = os.MkdirAll(dst, srcInfo.Mode())
  check("PRINT", err)

  // Open the directory and get all FileInfo, then initiate copying of
  // every source file/directory
  dir, _ := os.Open(src)
  objects, err := dir.Readdir(-1)
  for _, obj := range objects {
    srcFile, dstFile := src + "/" + obj.Name(), dst + "/" + obj.Name()

    if obj.IsDir() && obj.Name() != ".tmp" {
      // Recursively copy subdirectories while ignoring the .tmp directory
      copyDir(srcFile, dstFile)
    } else if !obj.IsDir() {
      // Simply copy the source file over since it's not a directory
      copyFile(srcFile, dstFile)
    }
  }
}

// Creates a snapshot of a given file/directory and outputs to the destination
func createSnapshot(src string, dst string) {
  if isDirectory(src) {
    copyDir(src, dst)
  } else {
    copyFile(src, dst)
  }
}