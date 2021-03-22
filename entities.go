package main

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

func (w Word) hasDefinition() bool {
	return len(w.Entries) > 0
}

type Cache map[string]Word
