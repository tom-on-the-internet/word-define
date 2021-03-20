package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	searchTerm, err := getSearchTerm()
	if err != nil {
		return err
	}

	config, err := getConfig()
	if err != nil {
		return err
	}

	if !config.valid() {
		return errInvalidConfig
	}

	word, err := getWord(searchTerm, config)
	if err != nil {
		return err
	}

	return printResult(word)
}

func getWord(searchTerm string, config Config) (Word, error) {
	if config.Cache {
		return getWordMaybeCached(searchTerm, config)
	}

	return getWordNotCached(searchTerm, config)
}

func getSearchTerm() (string, error) {
	flag.Parse()

	if len(flag.Args()) == 0 {
		return "", errNoSearchTerm
	}

	searchTerm := strings.ToLower(flag.Arg(0))

	return searchTerm, nil
}

func printResult(word Word) error {
	fmt.Println("[ " + strings.ToUpper(word.Spelling) + " ]")

	for index, entry := range word.Entries {
		num := strconv.Itoa(index + 1)

		fmt.Println()

		if len(word.Entries) > 1 {
			fmt.Println("(" + num + ")")
		}

		fmt.Println("DEFINITION: " + entry.Definition)

		if len(entry.Examples) > 0 {
			fmt.Println("EXAMPLES: " + strings.Join(entry.Examples, " | "))
		}

		if len(entry.Etymologies) > 0 {
			fmt.Println("ETYMOLOGIES: " + strings.Join(entry.Etymologies, " | "))
		}
	}

	return nil
}

func fetchWord(searchTerm string, config Config) (Word, error) {
	var (
		responseData OxfordResponse
		word         Word
	)

	url := "https://od-api.oxforddictionaries.com/api/v2/entries/en-gb/" + searchTerm + "?strictMatch=false"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("app_id", config.AppID)
	req.Header.Set("app_key", config.AppKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return word, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return word, err
	}

	return makeWordFromResponse(responseData)
}

func makeEntry(oxfordEntry OxfordEntry, oxfordSense OxfordSense) (Entry, error) {
	var entry Entry

	if len(oxfordSense.Definitions) == 0 {
		return entry, errNoDefinitionsFound
	}

	entry.Definition = oxfordSense.Definitions[0]
	entry.Etymologies = oxfordEntry.Etymologies

	for _, example := range oxfordEntry.Senses[0].Examples {
		entry.Examples = append(entry.Examples, example.Text)
	}

	return entry, nil
}

// Gets word from OxfordResponse.
// Unfortunately has to loop quite a bit.
func makeWordFromResponse(responseData OxfordResponse) (Word, error) {
	var word Word

	word.Spelling = responseData.Word

	results := responseData.Results

	for _, result := range results {
		for _, lexicalEntry := range result.LexicalEntries {
			for _, oxfordEntry := range lexicalEntry.Entries {
				for _, sense := range oxfordEntry.Senses {
					entry, err := makeEntry(oxfordEntry, sense)
					if err == nil {
						word.Entries = append(word.Entries, entry)
					}
				}
			}
		}
	}

	return word, nil
}

// Gets word from cache or remote source.
// This is run for users who have the cache enabled.
func getWordMaybeCached(searchTerm string, config Config) (Word, error) {
	var word Word

	cache, err := getCache()
	if err != nil {
		return word, err
	}

	word, ok := cache[searchTerm]
	if ok {
		if !word.hasDefinition() {
			return word, errNoDefinitionsFound
		}

		return word, nil
	}

	word, err = fetchWord(searchTerm, config)
	if err != nil {
		return word, err
	}

	cache[searchTerm] = word

	err = writeCache(cache)
	if err != nil {
		return word, err
	}

	if !word.hasDefinition() {
		return word, errNoDefinitionsFound
	}

	return word, nil
}

func getWordNotCached(searchTerm string, config Config) (Word, error) {
	word, err := fetchWord(searchTerm, config)
	if err != nil {
		return word, err
	}

	if word.Spelling == "" {
		return word, errNoDefinitionsFound
	}

	return word, nil
}
