package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	// Document containing all the words in English
	DictionaryUrl = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"

	// Document containing all links you need to scrape
	// Use: https://nx-public.s3-eu-west-1.amazonaws.com/Interview/endg-urls-short
	// for a quick dev feedback loop
	UrlListUrl = "https://nx-public.s3-eu-west-1.amazonaws.com/Interview/endg-urls"

	TextObtainers = 100
	LineSplitters = 100
	WordCounters  = 100
	BufferSize    = 100
)

func main() {
	dict, err := GetLines(DictionaryUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got dictionary of %d words\n", len(dict))

	wc := CreateWordCounter(dict)

	urls, err := GetLines(UrlListUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d urls to count\n", len(urls))

	wordCounts := wc.Count(urls)

	for word, count := range wordCounts {
		fmt.Printf("%s: %d times\n", word, count)
	}
}

type WordCounter struct {
	dictionary map[string]int
	counter    map[string]int
	locker     sync.Mutex

	urlPushersWG    sync.WaitGroup
	textObtainersWG sync.WaitGroup
	lineSplittersWG sync.WaitGroup
	wordCountersWG  sync.WaitGroup
}

func CreateWordCounter(dictionaryWords []string) *WordCounter {
	wc := WordCounter{
		dictionary: make(map[string]int),
		counter:    make(map[string]int),
	}

	for _, v := range dictionaryWords {
		wc.dictionary[v] = 0
	}

	return &wc
}

// Count counts how many times each word in dictionaryWords appears in all documents
// referenced in urlList. TODO(you): implement!
func (wc *WordCounter) Count(urlList []string) map[string]int {

	urlPipe := make(chan string, BufferSize)
	textLinePipe := make(chan string, BufferSize)
	wordPipe := make(chan string, BufferSize)

	wc.urlPushersWG.Add(1)
	go func() {
		var i int
		for _, url := range urlList {
			urlPipe <- url
			fmt.Println("Push URL #", i)
			i++
		}
		wc.urlPushersWG.Done()
	}()

	wc.textObtainersWG.Add(TextObtainers)
	for i := 0; i < TextObtainers; i++ {
		go wc.obtainText(urlPipe, textLinePipe)
	}

	wc.lineSplittersWG.Add(LineSplitters)
	for i := 0; i < LineSplitters; i++ {
		go wc.splitLineToWords(textLinePipe, wordPipe)
	}

	wc.wordCountersWG.Add(WordCounters)
	for i := 0; i < WordCounters; i++ {
		go wc.countWords(wordPipe)
	}

	wc.urlPushersWG.Wait()
	close(urlPipe)

	wc.textObtainersWG.Wait()
	close(textLinePipe)

	wc.lineSplittersWG.Wait()
	close(wordPipe)

	wc.wordCountersWG.Wait()

	return wc.counter
}

func (wc *WordCounter) countWords(wordPipe chan string) {
	counter := make(map[string]int)
	for word := range wordPipe {
		counter[word]++
	}

	wc.locker.Lock()
	for word, count := range counter {
		wc.counter[word] += count
	}
	wc.locker.Unlock()

	wc.wordCountersWG.Done()
}

func (wc *WordCounter) splitLineToWords(textLinePipe chan string, wordPipe chan string) {
	for line := range textLinePipe {
		words := strings.Split(line, " ")
		for _, word := range words {
			_, isExists := wc.dictionary[word]

			if isExists {
				wordPipe <- word
			}
		}
	}
	wc.lineSplittersWG.Done()
}

func (wc *WordCounter) obtainText(urlPipe chan string, textLinePipe chan string) {
	for url := range urlPipe {
		lines, err := GetLines(url)
		if err != nil {
			panic(err)
		}

		for _, line := range lines {
			textLinePipe <- line
		}
	}
	wc.textObtainersWG.Done()
}

// GetLines tries to http GET url and return the response body split by newline ("\n").
func GetLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
