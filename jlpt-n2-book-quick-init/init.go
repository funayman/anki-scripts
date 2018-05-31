package main

import (
  "database/sql"
  "encoding/json"
//  "errors"
//  "fmt"
//  "io"
  "log"
//  "os"
  "strconv"
  "strings"
//  "time"

  _ "github.com/mattn/go-sqlite3"
)

const (
  FieldsDelim = "\x1f"
  EnglishIndex = 0
  DisplayIndex = 1
  KanaIndex   = 2
  KanjiIndex  = 3
)

var (
	noteTypeNames = [...]string{
		"Japanese Vocab",
		"Japanese Vocab (and reversed card)",
	}
	noteTypeIDs  = make([]int64, 0)                // keep a map of note names and IDs
	models       = make(map[string]interface{}) // unmarshal JSON to map
	modelData []byte                            // raw JSON data from DB

  ankiDB *sql.DB
  bookDB *sql.DB
)

func init() {
  var err error
  //open Anki DB
  ankiDB, err = sql.Open("sqlite3", "./collection.anki2")
  if err != nil {
    log.Fatal(err)
  }

  bookDB, err = sql.Open("sqlite3", "./book.db")
  if err != nil {
    log.Fatal(err)
  }

  bookDB.Exec("DROP TABLE IF EXISTS word")
  bookDB.Exec("CREATE TABLE lindex(last_index int)")
  bookDB.Exec("CREATE TABLE word (id INT DEFAULT 0, eng VARCHAR, jpn VARCHAR, kana VARCHAR, kanji VARCHAR, bun_jpn VARCHAR, bun_eng VARCHAR, sound_word VARCHAR, sound_bun)")

  bookDB.Exec("DROP TABLE IF EXISTS lindex")
  bookDB.Exec("CREATE TABLE lindex (last_index)")
  bookDB.Exec("CREATE TABLE lindex (last_index INT)")
  bookDB.Exec("INSERT INTO lindex VALUES (1)")
}

func main() {
  // query data and unmarshal
  ankiDB.QueryRow("SELECT models FROM col").Scan(&modelData)
  json.Unmarshal(modelData, &models)

  // cycle through models and find the ones we're looking for
  for id, model := range models {
    //convert model to map to access data
    m, _ := model.(map[string]interface{})

    //check if current model matches one in noteTypeNames
    for _, name := range noteTypeNames {
      if name == m["name"] {
        mid, err := strconv.ParseInt(id, 10, 64)
        if err != nil {
          log.Fatal(err)
        } // if err
        noteTypeIDs = append(noteTypeIDs, mid)
      } // if names
    } // for _, name
  } // for id, model

  // query all cards for given note type
  rows, err := ankiDB.Query("SELECT id, flds FROM notes WHERE mid IN (?, ?)", noteTypeIDs[0], noteTypeIDs[1])
  if err != nil {
    log.Fatal(err)
  }

  // temp vars
  var id int64
  var flds string

  // cycle through each card
  for rows.Next() {
    err := rows.Scan(&id, &flds)
    if err != nil {
      log.Fatal(err)
    }

    // split the field data
    fields := strings.Split(flds, FieldsDelim)

    // grab all necessary field data
    eng := fields[EnglishIndex]
    dsp := fields[DisplayIndex]
    kna := fields[KanaIndex]
    knj := fields[KanjiIndex]

    addToBookDB(eng,dsp,kna,knj)
  }
}

func addToBookDB(english, display, kana, kanji string) {
  _, err := bookDB.Exec("INSERT INTO word (eng, jpn, kana, kanji) VALUES (?, ?, ?, ?)", english, display, kana, kanji)
  if err != nil {
    log.Fatal(err)
  }
}
