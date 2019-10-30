# Newspaper

Read webpages in readability mode, inside your terminal.

[![asciicast](https://asciinema.org/a/128518.png)](https://asciinema.org/a/128518)

### Why?
[Newsbeuter](http://newsbeuter.org/) is a great command line tool to read your favourite RSS feeds.
Also, it lets you choose which command to invoke when opening an article URL link, by setting the `browser` key inside configurations.
`newspaper` aims to be a simple command line tool to read URL's content in a clean and readable way.

You can choose to plug `newspaper` inside newsbeuter, or use it directly from the command line.

The heavy lifting is made by [go-readability](https://github.com/go-shiori/go-readability), the library that converts URL to markdown.

This package starts as a light and pluggable command between their api and the `less` command for Unix and `more` command for Windows.

### Usage

- `go get github.com/desmondhume/newspaper`
- `newspaper URL`


To use `newspaper` as newsbeuter browser, place this line in your newsbeuter config:

```
browser newspaper [OPTIONS] %u
```



### Options

```
-no-links    Remove markdown links
-plaintext   Disable ANSI codes (plaintext output)
-save-to-file    Save output to markdown file
```

### Todo

- [ ] Tests
- [x] Replace Mercury with a readability library
- [x] Save article to Markdown

