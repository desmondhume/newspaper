package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/logrusorgru/aurora"
	"github.com/lunny/html2md"
	"github.com/mitchellh/go-wordwrap"
)

func main() {
	var (
		nolinks    = flag.Bool("no-links", false, "Remove links")
		plaintext  = flag.Bool("plaintext", false, "Disable ANSI (plain-text output)")
		saveToFile = flag.Bool("save-to-file", false, "Save output to file")
	)

	flag.Parse()

	// Fetch article from given url
	articleURL := os.Args[len(os.Args)-1]
	parsedURL, _ := url.Parse(articleURL)

	article, err := readability.FromURL(parsedURL.String(), 5*time.Second)
	if err != nil {
		fmt.Printf("Unable to request data from URL %s: %v", articleURL, err)
	}

	// Convert html to readable markdown
	md := html2md.Convert(article.Content)
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
		output = regex.ReplaceAllString(output, fmt.Sprintf("%s", aurora.Bold("$1")))

		// Convert markdown wrappers to ANSI codes (to enhance subtitles)
		regex = regexp.MustCompile("## (.*)")
		output = regex.ReplaceAllString(output, fmt.Sprintf("%s", aurora.Bold("$1")))
	}

	// Wrap text to 80 columns to make the content more readable
	output = wordwrap.WrapString(output, 80)

	// Format article output with title and content
	output = fmt.Sprintf("%s\n%s", aurora.Bold(aurora.Red(article.Title)), output)

	if *saveToFile {
		outputAsBytes := []byte(output)
		filename := fmt.Sprintf("%s.md", article.Title)
		err = ioutil.WriteFile(filename, outputAsBytes, 0644)
		if err != nil {
			fmt.Printf("Unable to save output to the file: %v\n", err)
			os.Exit(1)
		}
	} else {
		cmd := exec.Command(PathToTerminalPagerProgram, ParamsForTerminalPagerProgram)

		// Set `less` stdin to string Reader
		cmd.Stdin = strings.NewReader(output)

		// Set `less` stdout to os stdout
		cmd.Stdout = os.Stdout

		// Start the command and wait for user actions
		err = cmd.Start()
		if err != nil {
			fmt.Print(err)
		} else {
			cmd.Wait()
		}
	}
}
