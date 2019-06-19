# linux-academy-to-anki
convert flashcards from [Linux Academy](https://linuxacademy.com) into a format that [Anki](https://apps.ankiweb.net/) can import.
The output best suited for Anki's *Basic* and *Basic (optional reverse card)*

## Usage
```
$ go run main.go deck1.json deck2.json
```

if you want a reverse card deck
```
$ go run main.go --with-reverse deck1.json deck2.json
```

## Getting Flash Card JSON Data
easiest to grab are the Instructor decks, which are usually pretty good.
If you want a user made made one, you'll have to do a bit more work

### Instructor Deck
```
# replace {:id} with the course id (check the URL)
https://linuxacademy.com/cp/flashcards/getInstructorDeck/module_id/{:id}
```

### User Submitted Deck
You'll need to find the index (usually starting from 2) of the deck you want to grab.
You can peek at the source code or just count if the list is small.
Run the following in the JS Console in your browser:

```js
(function(index){
  studyDeck = this.deck.decks[index]
  name = studyDeck.name;
  id = studyDeck.user_id;
  return $.ajax({
    url: "/cp/flashcards/getSharedDeck",
    dataType: "json",
    data: {"deck_name":name, "user_id":id},
    type: "post",
  }).success(function(data) {
    console.log(JSON.stringify(data).replace("\ufeff", ""));
  });
})(2)
```
Save the output to a JSON file.

### Your Own Decks
// TODO when I decide to make my own deck

## Importing to Anki
- Select `File` -> `Import...`
- Select your txt file to import (i.e. `anki_deck1.txt`)
- Be sure to set `Type` to `Basic` or `Basic (optional reverse card)`
- Choose your `deck`
- **IMPORTANT**: Make sure the your import options say: `Fields separated by: Semicolon`
- Go ahead and import
