# CS396 - Independent Study - Prof. Tim Richards

### What Should Be Included

  - `bigbrother.go`
  - `config.toml`
  - `csv.go`
  - `db.go`
  - `diff.go`
  - `json.go`
  - `main.go`
  - `README.md`
  - `snapshot.go`
  - `utils.go`

### Brief Overview

`bigbrother.go` contains the code that utilizes the
[fsnotify](https://fsnotify.org/) package to watch for any create, modify, or
remove operations occurring in a specific directory (and all of its
subdirectories). Diffs are generated upon each create, modify, or remove
operation, which are saved in a temporary directory and, if connected to the
internet, to a MongoDB database.

`csv.go` is, for the most part, not finished/very unpolished. It was written
for functionality at first, but later held off when JSON was the preferred
data format. As of right now, it isn't used in the watcher program at all.

`db.go` allows for connecting to a MongoDB database to store diffs generated
by the watcher.

`diff.go` provides the watcher program the ability to generate diffs between two
files.

`json.go` contains the collection of fields used to contain each diff, and
the collection of all diffs that are stored locally and on a MongoDB database.

`main.go` simply serves as the starting point of the watcher program where it
creates a temporary directory (or recreates if it already exist) for storing
diffs, reads configuration data set in `config.toml`, creates a snapshot of the
current directory, and starts up the watcher.

`snapshot.go` allows for the watcher to be able to copy files and directories
over to specified directory, creating a snapshot of everything at the moment.

`utils.go` is mainly for the purpose of code refactoring and keeping everything
clean. It contains various functions that are called throughout the watcher.

`config.toml` simply contains user-defined configuration settings that are used
in the watcher program.

### Usage

This assumes you already have the GOPATH environment variable set. If you don't
think you do, follow [this](https://github.com/golang/go/wiki/GOPATH).

First, go to the directory:

```bash
$ cd path/to/directory
```

Packages will need to be installed before running the watcher program.

If on Linux/Unix
```bash
$ chmod 777 packages.sh
$ ./packages.sh
```

Otherwise, run:
```bash
$ go get github.com/fsnotify/fsnotify
$ go get gopkg.in/mgo.v2
$ go get github.com/pmezard/go-difflib/difflib
$ go get github.com/BurntSushi/toml
```

Now you should be able to run the watcher program:
```bash
$ go run *.go
```

If you would like to build the watcher program into a binary:
```bash
$ go build
```

A new file should appear called ```cs396```, which you can run using:
```bash
$ ./cs396
```
