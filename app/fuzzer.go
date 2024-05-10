package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"log"
)

func main() {

	var (
		flURL      = flag.String("url", "", "URL to fuzz")
		flWordlist = flag.String("wordlist", "", "Wordlist to use")
	)

	flag.Parse()

	if *flURL == "" {
		fmt.Println("URL is required: --url")
		return
	}

	if *flWordlist == "" {
		fmt.Println("Wordlist is required: --wordlist")
		return
	}

	url := strings.TrimSpace(*flURL)

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	fmt.Println("Fuzzing: ", url)


	wordlistPath := strings.TrimSpace(*flWordlist)
	wordlist := readWordlist(wordlistPath)

	fuzzUrl(url, wordlist)
}

func fuzzUrl(url string, wordlist []string) {
	for _, word := range wordlist {
		newUrl := url + word
		go checkUrl(newUrl)
	}
}

func checkUrl(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("URL exists: ", url)
	}
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
