package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	errNoDefinitionsFound = errors.New("no definitions found")
	errNoSearchTerm       = errors.New("no search term provided")
	errInvalidConfig      = errors.New("config invalid. must have a valid app key and id")
)

type Config struct {
	AppKey string `json:"appKey"`
	AppID  string `json:"appId"`
	Cache  bool   `json:"cache"`
}

func (c Config) valid() bool {
	return c.AppKey != "" && c.AppID != ""
}

type Entry struct {
	Definition  string   `json:"definition"`
	Examples    []string `json:"examples"`
	Etymologies []string `json:"etymologies"`
}

type Word struct {
	Spelling string
	Entries  []Entry `json:"entries"`
}

type OxfordSense struct {
	Definitions   []string `json:"definitions"`
	DomainClasses []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"domainClasses"`
	Examples []struct {
		Text      string `json:"text"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers,omitempty"`
	} `json:"examples"`
	ID              string `json:"id"`
	SemanticClasses []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"semanticClasses"`
	ShortDefinitions []string `json:"shortDefinitions"`
	Registers        []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"registers,omitempty"`
	Subsenses []struct {
		Definitions   []string `json:"definitions"`
		DomainClasses []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domainClasses"`
		Examples []struct {
			Text string `json:"text"`
		} `json:"examples"`
		ID              string `json:"id"`
		SemanticClasses []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"semanticClasses"`
		ShortDefinitions []string `json:"shortDefinitions"`
	} `json:"subsenses,omitempty"`
	Synonyms []struct {
		Language string `json:"language"`
		Text     string `json:"text"`
	} `json:"synonyms,omitempty"`
	ThesaurusLinks []struct {
		EntryID string `json:"entry_id"`
		SenseID string `json:"sense_id"`
	} `json:"thesaurusLinks,omitempty"`
}

type OxfordEntry struct {
	Etymologies     []string `json:"etymologies"`
	HomographNumber string   `json:"homographNumber"`
	Pronunciations  []struct {
		AudioFile        string   `json:"audioFile"`
		Dialects         []string `json:"dialects"`
		PhoneticNotation string   `json:"phoneticNotation"`
		PhoneticSpelling string   `json:"phoneticSpelling"`
	} `json:"pronunciations"`
	Senses []OxfordSense `json:"senses"`
}

type OxfordResponse struct {
	ID       string `json:"id"`
	Metadata struct {
		Operation string `json:"operation"`
		Provider  string `json:"provider"`
		Schema    string `json:"schema"`
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Entries         []OxfordEntry `json:"entries"`
			Language        string        `json:"language"`
			LexicalCategory struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Phrases []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"phrases"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
	Word string `json:"word"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var word Word

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

	if config.Cache {
		word, err = getWordMaybeCached(searchTerm, config)
		if err != nil {
			return err
		}
	} else {
		word, err = getWordNotCached(searchTerm, config)
		if err != nil {
			return err
		}
	}

	return printResult(word)
}

func getSearchTerm() (string, error) {
	flag.Parse()

	if len(flag.Args()) == 0 {
		return "", errNoSearchTerm
	}

	searchTerm := strings.ToLower(flag.Arg(0))

	return searchTerm, nil
}

func getConfig() (Config, error) {
	var config Config

	err := createConfigIfNotExists()
	if err != nil {
		return config, err
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}

	appConfigDir := userConfigDir + "/word-define"
	configFilename := appConfigDir + "/config.json"

	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		// json was invalid, write empty json
		configJSON, _ := json.MarshalIndent(config, "", "  ")
		err = ioutil.WriteFile(configFilename, configJSON, 0600)
	}

	return config, err
}

func getCache() (map[string]Word, error) {
	wordMap := make(map[string]Word)

	err := createCacheIfNotExists()
	if err != nil {
		return wordMap, err
	}

	cacheFilename, err := getCacheFilename()
	if err != nil {
		return wordMap, err
	}

	bytes, err := ioutil.ReadFile(cacheFilename)
	if err != nil {
		return wordMap, err
	}

	err = json.Unmarshal(bytes, &wordMap)
	if err != nil {
		// json was invalid, write empty json
		ioutil.WriteFile(cacheFilename, []byte("{}"), 600)
	}

	return wordMap, nil
}

func createCacheIfNotExists() error {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	appCachedir := userCacheDir + "/word-define"
	cacheFilename := appCachedir + "/dict.json"

	// create directory if not exists
	_ = os.Mkdir(appCachedir, 0700)

	// create file if not exists
	if _, err := os.Stat(cacheFilename); os.IsNotExist(err) {
		file, err := os.OpenFile(cacheFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
		if err != nil {
			return err
		}

		file.Close()

		return ioutil.WriteFile(cacheFilename, []byte("{}"), 0600)
	}

	return nil
}

func createConfigIfNotExists() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appConfigdir := userConfigDir + "/word-define"
	configFilename := appConfigdir + "/config.json"

	// create directory if not exists
	_ = os.Mkdir(appConfigdir, 0700)

	// create file if not exists
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		file, err := os.OpenFile(configFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
		if err != nil {
			return err
		}

		file.Close()

		var config Config

		configJSON, _ := json.MarshalIndent(config, "", "  ")

		return ioutil.WriteFile(configFilename, configJSON, 0600)
	}

	return nil
}

func getCacheFilename() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	appCachedir := userCacheDir + "/word-define"
	cacheFilename := appCachedir + "/dict.json"

	return cacheFilename, err
}

func writeCache(wordMap map[string]Word) error {
	cacheFilename, err := getCacheFilename()
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(cacheFilename, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	wordMapJSON, _ := json.MarshalIndent(wordMap, "", "  ")

	return ioutil.WriteFile(cacheFilename, wordMapJSON, 0600)
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
		return entry, errors.New("No definition.")
	}

	entry.Definition = oxfordSense.Definitions[0]
	entry.Etymologies = oxfordEntry.Etymologies

	for _, example := range oxfordEntry.Senses[0].Examples {
		entry.Examples = append(entry.Examples, example.Text)
	}

	return entry, nil
}

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
func getWordMaybeCached(searchTerm string, config Config) (Word, error) {
	var word Word

	wordMap, err := getCache()
	if err != nil {
		return word, err
	}

	word, ok := wordMap[searchTerm]
	if ok {
		if word.Spelling == "" {
			return word, errNoDefinitionsFound
		}

		return word, nil
	}

	word, err = fetchWord(searchTerm, config)
	if err != nil {
		return word, err
	}

	wordMap[searchTerm] = word

	err = writeCache(wordMap)
	if err != nil {
		return word, err
	}

	if word.Spelling == "" {
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
