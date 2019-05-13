package main

import (
	"testing"
)

func TestSanityWordCounter(t *testing.T) {
	mockDictionary := "https://s3-us-west-1.amazonaws.com/nx-rotem/wc/words.txt"
	mockUrls := "https://s3-us-west-1.amazonaws.com/nx-rotem/wc/urls.txt"

	words, err := GetLines(mockDictionary)
	if err != nil {
		t.Fatal(err)
	}

	urls, err := GetLines(mockUrls)
	if err != nil {
		t.Fatal(err)
	}

	wc := &WordCounter{}
	count, err := wc.Count(urls, words)

	testCases := []struct {
		word          string
		expectedCount int
	}{
		// word that exists in the dict and an item
		{word: "hello", expectedCount: 3},

		// exists in the dict but not in an item
		{word: "world", expectedCount: 0},

		// exists in both
		{word: "goodbye", expectedCount: 1},

		// exists only in the item
		{word: "fox", expectedCount: 0},
	}

	for _, c := range testCases {
		t.Run(c.word, func(t *testing.T) {
			actual := count[c.word]
			if actual != c.expectedCount {
				t.Fatalf("%s count should be %d", c.word, c.expectedCount)
			}
		})
	}
}
