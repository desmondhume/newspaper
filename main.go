package main

import (
	"flag"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"github.com/lunny/html2md"
	"github.com/mitchellh/go-wordwrap"
	"html"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"github.com/go-shiori/go-readability"
	"time"
)

func main() {
	var (
		nolinks   = flag.Bool("no-links", false, "Remove links")
		plaintext = flag.Bool("plaintext", false, "Disable ANSI (plain-text output)")
	)

	flag.Parse()

	// Fetch article from given url
	articleUrl := os.Args[len(os.Args)-1]
	parsedURL, _ := url.Parse(articleUrl)

	article, err := readability.FromURL(parsedURL, 5*time.Second)
	if err != nil {
		fmt.Printf("Unable to request data from URL %s: %v", articleUrl, err)
	}

	// Convert html to readable markdown
	md := html2md.Convert(article.RawContent)
	output := html.UnescapeString(md)

	var regex *regexp.Regexp

	// Squash multiple lines blocks into single blank lines
	regex = regexp.MustCompile(`(\s*\n){2,}`)
	output = regex.ReplaceAllString(output, "\n\n")

	// Remove leading whitespaces
	regex = regexp.MustCompile(`[\n\n ][ \t]+`)
	output = regex.ReplaceAllString(output, "")

	// Remove links if -nolinks is passed
	if *nolinks {
		// Remove empty links (like js-driven anchors)
		regex = regexp.MustCompile(`\[\]\(\)`)
		output = regex.ReplaceAllString(output, "")

		// Remove other links
		regex = regexp.MustCompile(`\[(.*)\]\((.*)\)`)
		output = regex.ReplaceAllString(output, "$1")
	}

	if !*plaintext {
		// Convert markdown wrappers to ANSI codes (to enhance subtitles)
		regex = regexp.MustCompile(`\*\*(.*)\*\*`)
		output = regex.ReplaceAllString(output, fmt.Sprintf("%s", Bold("$1")))

		// Convert markdown wrappers to ANSI codes (to enhance subtitles)
		regex = regexp.MustCompile("## (.*)")
		output = regex.ReplaceAllString(output, fmt.Sprintf("%s", Bold("$1")))
	}

	// Wrap text to 80 columns to make the content more readable
	output = wordwrap.WrapString(output, 80)

	// Format article output with title and content
	output = fmt.Sprintf("%s\n%s", Bold(Red(article.Meta.Title)), output)
	cmd := exec.Command("/usr/bin/less", "-s")

	// Set `less` stdin to string Reader
	cmd.Stdin = strings.NewReader(output)

	// Set `less` stdout to os stdout
	cmd.Stdout = os.Stdout

	// Start the command and wait for user actions
	cmd.Start()
	cmd.Wait()
}
