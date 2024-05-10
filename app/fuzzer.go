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

	fmt.Println("Fuzzing: ", *flURL)

	url := strings.TrimSpace(*flURL)

	wordlistPath := strings.TrimSpace(*flWordlist)
	wordlist := readWordlist(wordlistPath)

	fuzzUrl(url, wordlist)
}

func fuzzUrl(url string, wordlist []string) {
	for _, word := range wordlist {
		newUrl := fmt.Sprintf("%s/%s", url, word)
		exists := checkUrl(newUrl)
		if exists {
			fmt.Println("URL exists: ", newUrl)
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
