package main

import (
  "fmt"
  "os"
  "log"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
)

const (
  SoundFormat = "[sound:%s]"
)

var (
  bookDB *sql.DB
)

func addWordsToDb(words []string) {
  for _, word := range words {
    fmt.Println(word)
  }
}

func getLastIndex() (index int) {
  bookDB.QueryRow("SELECT * FROM lindex").Scan(&index)
  return
}

func usageAndExit() {
  fmt.Print("usage: add [command]\n\twords\twords to be added (must have at least one word)\n\tindex\treturns the starting index of the word to be added\n")
  os.Exit(1)
}

func init() {
  var err error
  //open Anki DB
  bookDB, err = sql.Open("sqlite3", "./book.db")
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  args := os.Args
  if len(args) < 2 {
    usageAndExit()
  }

  if command := args[1]; command == "words" {
    addWordsToDb(args[2:])
  } else if command == "index" {
    fmt.Printf("Current Index in DB: %d\n", getLastIndex())
  } else {
    usageAndExit()
  }
}
