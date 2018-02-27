package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	FieldsDelim = "\x1f"
	AudioIndex  = 5 // customize based on your field index
	KanjiIndex  = 3 // customize based on your field index
	KanaIndex   = 2 // customize based on your field index
)

var (
	noteTypeNames = [...]string{
		"Japanese Vocab",
		"Japanese Vocab (and reversed card)",
	} // customize for your note types
	noteTypes = make(map[string]int64)       // keep a map of note names and IDs
	models    = make(map[string]interface{}) // unmarshal JSON to map
	modelData []byte                         // raw JSON data from DB
	cards     []Card                         // hold our card types to update the DB after updating
)

type Card struct {
	Fields string
	Id     int64
}

func downloadAudio(kanji, kana string) error {
	BaseURLF := "http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kanji=%s&kana=%s"

	fullURL := fmt.Sprintf(BaseURLF, kanji, kana)

	//Do the downloading
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}

	if resp.StatusCode == 500 {
		return errors.New("Sorry! Internal Server error! Try again in a bit")
	}

	//The audio for this clip is currently unavailable
	if resp.Header.Get("Content-Length") == "52288" {
		return errors.New(fmt.Sprintf("Couln't find audio for %s[%s] -- check kanji/kana", kanji, kana))
	}

	fileName := fmt.Sprintf("%s_%s.mp3", kanji, kana)
	audioFile, err := os.Create(fileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating file %s[%s]: %s", kanji, kana, err))
	}

	defer audioFile.Close()
	defer resp.Body.Close()

	//Write the body to the file
	_, err = io.Copy(audioFile, resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error saving file %s[%s]: %s", kanji, kana, err))
	}

	return nil
}

func truncate(kanji, kana string, size int) (string, string) {
	tmpKanjiRune := []rune(kanji)
	tmpKanaRune := []rune(kana)

	kanji = string(tmpKanjiRune[:len(tmpKanjiRune)-size])
	kana = string(tmpKanaRune[:len(tmpKanaRune)-size])

	return kanji, kana
}

func main() {
	//open DB
	db, err := sql.Open("sqlite3", "./collection.anki2")
	if err != nil {
		log.Fatal(err)
	}

	// query data and unmarshal
	db.QueryRow("SELECT models FROM col").Scan(&modelData)
	json.Unmarshal(modelData, &models)

	// cycle through models and find the ones we're looking for
	for id, model := range models {
		//convert model to map to access data
		m, _ := model.(map[string]interface{})

		//check if current model matches one in noteTypeNames
		for _, name := range noteTypeNames {
			if name == m["name"] {
				noteTypes[name], err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					log.Fatal(err)
				} // if err
			} // if names
		} // for _, name
	} // for id, model

	log.Printf("Current noteTypes %v\n", noteTypes)

	// go through each note type and process cards
	for name, mid := range noteTypes {
		log.Printf("Processing cards for %s (ID %d)\n", name, mid)

		// query all cards for given note type
		rows, err := db.Query("SELECT id, flds FROM notes WHERE mid=?", mid)
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
			kanji := fields[KanjiIndex]
			kana := fields[KanaIndex]
			audio := fields[AudioIndex]

			// already have audio, skip
			if audio != "" {
				continue
			}

			// no kanji for card (possibly katakana or word has no kanji)
			if kanji == "" {
				kanji = kana
			}

			// my cards have a する suffix for suru verbs and a な suffix for na-adjectives
			// these cause a "audio not found" on JapanesePod101, they need to be removed
			// in order to download the audio correctly
			tmpRune := []rune(kanji)

			if len(tmpRune) >= 3 && tmpRune[len(tmpRune)-1] == 'る' && tmpRune[len(tmpRune)-2] == 'す' {
				kanji, kana = truncate(kanji, kana, 2)
			} else if len(tmpRune) >= 2 && tmpRune[len(tmpRune)-1] == 'な' {
				kanji, kana = truncate(kanji, kana, 1)
			}

			// we good, lets do this shit!
			// download audio
			err = downloadAudio(kanji, kana)
			if err != nil {
				// error downloading file
				// continue marching forward and address afterwards
				fmt.Println(err)
				continue
			}

			// replace the audio field with an Anki sound URI
			fields[AudioIndex] = fmt.Sprintf("[sound:%s_%s.mp3]", kanji, kana)

			// join our fields back together
			// and assign data to our card type
			// so we can update the DB
			newFields := strings.Join(fields, FieldsDelim)
			cards = append(cards, Card{Fields: newFields, Id: id})

			//dont overload the server now...
			time.Sleep(time.Millisecond * 1250)
		}
	}

	// update the DB
	for _, card := range cards {
		stmt, err := db.Prepare("UPDATE notes SET flds=?, mod=?, usn=? WHERE id=?")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(card.Fields, time.Now().Unix(), -1, card.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
