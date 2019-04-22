package models

import (
	"strings"
)

// GeneratePuns given query word, homophone word list, and phrase
func GeneratePuns(word string, homophoneWords []WordRow, phrases []Phrase) []string {
	puns := []string{}

	// TODO: remove punctuation
	for _, phrase := range phrases {
		tokens := strings.Split(phrase.PhraseText, " ")

		for i, token := range tokens {
			for _, homophoneWord := range homophoneWords {
				if strings.ToLower(token) == homophoneWord.Word {
					tokens[i] = word
					break
				}
			}
		}
		result := strings.Join(tokens, " ")
		puns = append(puns, result)
	}

	return puns
}
