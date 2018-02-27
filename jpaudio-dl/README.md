# jpaudio-dl
A small command line interface to download Japanese audio from JapanesePod101.

You can read the blog post about it at https://funayman.me/posts/anki-hacking-with-go/

## Usage
```bash
# word with both kanji and kana
$ jpaudio-dl 漢字 かんじ
File 漢字_かんじ.mp3 downloaded!

# word without kanji equivalent
$ jpaudio-dl もう もう
File もう_もう.mp3 downloaded!

# word in katakana
$ jpaudio-dl ナビ ナビ
File ナビ_ナビ.mp3 downloaded!
```
