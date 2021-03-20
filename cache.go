package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func getCache() (Cache, error) {
	var wordMap Cache

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

	// create directory if not exists
	_ = os.Mkdir(userCacheDir, 0700)

	appCachedir := userCacheDir + "/word-define"
	cacheFilename := appCachedir + "/dict.json"

	// create directory if not exists
	_ = os.Mkdir(appCachedir, 0700)

	// create file if not exists
	if _, err := os.Stat(cacheFilename); os.IsNotExist(err) {
		return ioutil.WriteFile(cacheFilename, []byte("{}"), 0600)
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
