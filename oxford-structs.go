package main

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
