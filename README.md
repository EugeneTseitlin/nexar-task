Website Vocabulary Analyzer Test
================================

Hello! Welcome!

In this test you are required to write a java program that counts the number of occurrences of certain words in a group of websites.


Input
=====
* 'vocabulary' - A vocabulary of predefined words. Those are the only words that we are interested in counting.

* 'urls' - A list of URLs to html pages.


Goal
====
Fetch each of the pages from the web, count the number of occurrences of each word from the dictionary in them, and print out the results.


Gotchas
=======
* We are only interested in full-word match. No need to look for partial or fuzzy matches.

* To make things easier, you do not have to distinguish between the html structure and body. You can treat the entire html as a text document.

* Your code should be clean and well organized.

* The list of urls is fairly long, which could cause the program to run for a long time. Eventually, we'd like for you to try and make it run as fast as possible.

* Consider that there is no limit on the size of the fetched URLs

* Consider that some websites may respond very slowly 


Code Structure
==============
* main.go - the entry point for the program which you will be running. Main already initializes the input params (vocabulary and urls).

in here you will find a struct `WordCounter` with a method `Count(urlList, dictionaryWords []string) (map[string]int, error)` which you need to implement.

```bash
# running the code
➜  word-counter-task-go go run main.go
got dictionary of 466551 words
got 105 urls to count

... more output pending your work

```

* word_counter_test.go - to make development easier we've written a sanity test.

```bash
# running the test
➜  word-counter-task-go go test .
--- FAIL: TestSanityWordCounter (2.83s)
    --- FAIL: TestSanityWordCounter/hello (0.00s)
        word_counter_test.go:37: hello count should be 3
    --- FAIL: TestSanityWordCounter/goodbye (0.00s)
        word_counter_test.go:37: goodbye count should be 1
FAIL
FAIL    github.com/rotemtam/word-counter-task-go        2.853s

# will PASS after you implement your code!
```
