package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	BaseURLF = "http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kanji=%s&kana=%s"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatal("usage: jpodaudio 漢字 かな")
	}

	kanji := args[1]
	kana := args[2]
	fullUrl := fmt.Sprintf(BaseURLF, kanji, kana)

	//Do the downloading
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 500 {
		log.Fatal("Sorry! Internal Server error! Try again in a bit")
	}

	//The audio for this clip is currently unavailable
	if resp.Header.Get("Content-Length") == "52288" {
		log.Fatal("Sorry! Audio doesn't exist! Check kanji and kana or JapanesePod101 doesn't actually have the audio")
	}

	fileName := fmt.Sprintf("%s_%s.mp3", kanji, kana)
	audioFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer audioFile.Close()
	defer resp.Body.Close()

	//Write the body to the file
	_, err = io.Copy(audioFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File %s.mp3 downloaded!\n", fileName)
}
