package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Cards []struct {
	CreatedAt    float64 `json:"created_at"`
	Likes        []int   `json:"likes,omitempty"`
	Dislikes     []int   `json:"dislikes,omitempty"`
	DeckName     string  `json:"deck_name"`
	SpamCount    int     `json:"spam_count,omitempty"`
	UID          string  `json:"uid"`
	Front        string  `json:"front"`
	SpamVotes    []int   `json:"spam_votes,omitempty"`
	UpdatedAt    float64 `json:"updated_at"`
	ForksCount   int     `json:"forks_count,omitempty"`
	UserID       int     `json:"user_id"`
	Back         string  `json:"back"`
	LikesBalance int     `json:"likes_balance,omitempty"`
	Puid         string  `json:"puid,omitempty"`
}

func (c Cards) String() string {
	sb := strings.Builder{}
	for i := 0; i < len(c); i++ {
		card := c[i]
		newFront := strings.Replace(card.Front, `"`, `""`, -1)
		newBack := strings.Replace(card.Back, `"`, `""`, -1)
		if withReverse {
			sb.WriteString(fmt.Sprintf("\"%s\"; \"%s\"; Y\n", newFront, newBack))
		} else {
			sb.WriteString(fmt.Sprintf("\"%s\"; \"%s\"\n", newFront, newBack))
		}
	}
	return sb.String()
}

var (
	withReverse bool
)

func init() {
	flag.BoolVar(&withReverse, "with-reverse", false, "changes note type to Basic (optional reversed card)")
	flag.Parse()
}

func main() {
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("at least one file name is required")
		os.Exit(1)
	}

	for _, file := range files {
		_, baseName := filepath.Split(file)
		outputFilename := fmt.Sprintf("anki_%s.txt", baseName[:len(baseName)-len(filepath.Ext(baseName))])

		f, err := os.Open(file)
		exitOnErr(err)
		defer f.Close()

		input, err := ioutil.ReadAll(f)

		var data Cards
		err = json.Unmarshal(input, &data)
		exitOnErr(err)

		outf, err := os.Create(outputFilename)
		exitOnErr(err)
		defer outf.Close()

		outf.WriteString(data.String())
	}
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
