package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	// Arguments
	var (
		flURL      = flag.String("url", "", "URL to fuzz")
		flWordlist = flag.String("wordlist", "", "Wordlist to use")
		flOutput   = flag.String("output", "", "Output file")
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
	words := make(chan string, 100)
	results := make(chan string)
	var urls []string

	// Fuzzing
	for i := 0; i < cap(words); i++ {
		go fuzzUrl(url, words, results)
	}

	// Loading wordlist in channel
	go func() {
		for _, word := range wordlist {
			words <- word
		}
	}()

	// Collecting results
	for _, word := range wordlist {
		url := <-results
		if word != "" {
			urls = append(urls, url)
		}
	}

	close(words)
	close(results)

	// Write results
	output := strings.TrimSpace(*flOutput)

	if output != "" {
		writeResults(urls, output)
	}
}

func fuzzUrl(url string, words, results chan string) {
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
	}
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
