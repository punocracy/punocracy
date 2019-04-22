package models

import (
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

// GeneratePuns given query word, homophone word list, and phrase
func GeneratePuns(word string, homophoneWords []WordRow, phrases []Phrase) []string {
	puns := []string{}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		logrus.Error(err)
	}

	for _, phrase := range phrases {
		processedString := reg.ReplaceAllString(phrase.PhraseText, "")
		tokens := strings.Split(processedString, " ")

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
