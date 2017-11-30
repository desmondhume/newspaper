package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"github.com/lunny/html2md"
	"github.com/mitchellh/go-wordwrap"
	"html"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type FeedItem struct {
	Title   string
	Content string
}

func main() {
	var (
		nolinks   = flag.Bool("no-links", false, "Remove links")
		plaintext = flag.Bool("plaintext", false, "Disable ANSI (plain-text output)")
	)

	flag.Parse()
	item := FeedItem{}

	// Create client to fetch article from given url
	client := &http.Client{}
	url := fmt.Sprintf("https://mercury.postlight.com/parser?url=%s", url.QueryEscape(os.Args[len(os.Args)-1]))
	request, _ := http.NewRequest("GET", url, nil)

	// Append mercury postlight api key to request headers
	apikey := os.Getenv("MERCURY_API_KEY")
	if apikey == "" {
		fmt.Printf("API key not found. Set MERCURY_API_KEY in your terminal.")
		os.Exit(1)
	}
	request.Header.Set("x-api-key", apikey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Unable to request data from URL %s: %v", url, err)
	}
	defer response.Body.Close()

	// Check for messages from server
	if response.StatusCode != 200 {
		errMsg := struct {
			Message string `json:"message"`
		}{}
		err := json.NewDecoder(response.Body).Decode(&errMsg)
		if err != nil {
			fmt.Printf("Unable to decode error message from server: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Error: Received message from server: %s\n", errMsg.Message)
		os.Exit(1)
	}

	// Fill the `item` attributes with api response
	if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
		fmt.Printf("Unable to decode json from Mercury: %v\n", err)
		os.Exit(1)
	}

	// Convert html to readable markdown
	md := html2md.Convert(item.Content)
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
	output = fmt.Sprintf("%s\n%s", Bold(Red(item.Title)), output)
	cmd := exec.Command("/usr/bin/less", "-s")

	// Set `less` stdin to string Reader
	cmd.Stdin = strings.NewReader(output)

	// Set `less` stdout to os stdout
	cmd.Stdout = os.Stdout

	// Start the command and wait for user actions
	cmd.Start()
	cmd.Wait()
}
