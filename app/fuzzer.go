package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
  "github.com/schollz/progressbar/v3" 
)

type empty struct{}

func main() {

	// Arguments
	var (
		flURL         = flag.String("url", "", "URL to fuzz")
		flWorkerCount = flag.Int("t", 100, "Workers to use")
		flWordlist    = flag.String("wordlist", "", "Wordlist to use")
		flOutput      = flag.String("output", "", "Output file")
	)

	flag.Parse()

	// Mandatory arguments
	if *flURL == "" {
		fmt.Println("URL is required: --url")
		return
	}

	if *flWordlist == "" {
		fmt.Println("Wordlist is required: --wordlist")
		return
	}

	// URL parsing
	url := strings.TrimSpace(*flURL)

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	fmt.Println("Fuzzing: ", url)

	// Wordlist parsing
	wordlistPath := strings.TrimSpace(*flWordlist)
	wordlist := readWordlist(wordlistPath)

	// Create channels
	words := make(chan string, *flWorkerCount)
	results := make(chan string)
	progress := make(chan int)
	tracker := make(chan empty)

  // Progress
  bar := progressbar.Default(int64(len(wordlist)))
	processedWords := 0

	// Results
	var urls []string

	// Print progress
	go func() {
		defer close(progress)
		for range progress {
			processedWords++
      bar.Add(1)
		}
	}()

	// Fuzzing
	for i := 0; i < *flWorkerCount; i++ {
		go fuzzUrl(url, words, results, progress, tracker)
	}

	// Results
	go func() {
		for range wordlist {
			url := <-results
			if url != "" {
				urls = append(urls, url)
        updateScreen(urls)
			}
		}
		var e empty
		tracker <- e
	}()

	// Loading wordlist in channel
	for _, word := range wordlist {
		words <- word
	}

	// Close channels
	close(words)

	for i := 0; i < *flWorkerCount; i++ {
		<-tracker
	}

	close(results)
	<-tracker

	fmt.Println("Fuzzing finished!")

	// Write results
	output := strings.TrimSpace(*flOutput)

	if output != "" {
		fmt.Println("Writing results to: ", output)
		writeResults(urls, output)
	}
}

func fuzzUrl(url string, words, results chan string, progress chan int, tracker chan empty) {
	for word := range words {
		if word == "" {
			continue
		}
		newUrl := url + word
		exists := checkUrl(newUrl)
		if exists {
			fmt.Println("URL exists: ", newUrl)
			results <- newUrl
		} else {
			results <- ""
		}
		progress <- 1
	}
	var e empty
	tracker <- e
}

func checkUrl(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK
}

func readWordlist(wordlistPath string) []string {
	readFile, err := os.Open(wordlistPath)

	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	var wordlist []string

	for fileScanner.Scan() {
		if !strings.HasPrefix(fileScanner.Text(), "#") {
			wordlist = append(wordlist, fileScanner.Text())
		}
	}

	return wordlist
}

func writeResults(urls []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for _, url := range urls {
		file.WriteString(url + "\n")
	}
}

// Screen

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printResults(urls []string) {
  fmt.Println("Urls Discovered:")
	for _, url := range urls {
		fmt.Println(url)
	}
  fmt.Println("")
}

func updateScreen(urls []string) {
	clearScreen()
	printResults(urls)
}
