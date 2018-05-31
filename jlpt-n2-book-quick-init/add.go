package main

import (
  "fmt"
  "os"
  "log"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
)

const (
  SoundWordFormat = "[sound:N2v_%04d.mp3]"
  SoundBunFormat = "[sound:N2v_%04ds.mp3]"
)

var (
  bookDB *sql.DB
)

func addWordsToDb(words []string) {
  for _, word := range words {
    index := getLastIndex()
    var rslt string
    bookDB.QueryRow("SELECT kanji FROM word WHERE kanji = ?", word).Scan(&rslt)

    if rslt == "" {
      bookDB.Exec("INSERT INTO word (id, kanji, sound_word, sound_bun) VALUES (?, ?, ?, ?)",
        index, word, fmt.Sprintf(SoundWordFormat, index), fmt.Sprintf(SoundBunFormat, index))
    } else {
      bookDB.Exec("UPDATE word SET id = ?, sound_word = ?, sound_bun = ? WHERE kanji = ?",
        index, fmt.Sprintf(SoundWordFormat, index), fmt.Sprintf(SoundBunFormat, index), word)
    }

    // update the index
    bookDB.Exec("UPDATE lindex SET last_index = ?", index + 1)
  }
}

func getLastIndex() (index int) {
  bookDB.QueryRow("SELECT * FROM lindex").Scan(&index)
  return
}

func exportAll() {
  rows, err := bookDB.Query("SELECT * FROM word WHERE id IS NOT 0 ORDER BY id")
  if err != nil {
    log.Fatal(err)
  }

  var id, eng, jpn, kana, kanji, bun_jpn, bun_eng, sound_word, sound_bun sql.NullString
  for rows.Next() {
    err := rows.Scan(&id, &eng, &jpn, &kana, &kanji, &bun_jpn, &bun_eng, &sound_word, &sound_bun)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Printf("%s|%s|%s|%s|%s|%s|%s|%s|%s\n", id.String, eng.String, jpn.String, kana.String, kanji.String, bun_jpn.String, bun_eng.String, sound_word.String, sound_bun.String)
  }
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
  } else if command == "export" {
    exportAll()
  } else {
    usageAndExit()
  }
}
