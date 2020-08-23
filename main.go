package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DictionaryUrl = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"
	UrlListUrl    = "https://nx-public.s3-eu-west-1.amazonaws.com/Interview/endg-urls"
)

func main() {
	wc := &WordCounter{}
	dict, err := GetLines(DictionaryUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got dictionary of %d words\n", len(dict))

	urls, err := GetLines(UrlListUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d urls to count\n", len(urls))

	wordCounts, err := wc.Count(urls, dict)
	if err != nil {
		panic(err)
	}

	for word, count := range wordCounts {
		fmt.Printf("%s: %d times\n", word, count)
	}
}

type WordCounter struct {
}

// Count counts how many times each word in dictionaryWords appears in all documents
// referenced in urlList. TODO(you): implement!
func (*WordCounter) Count(urlList, dictionaryWords []string) (map[string]int, error) {
	return nil, nil
}

// GetLines tries to http GET url and return the response body split by newline ("\n").
func GetLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
