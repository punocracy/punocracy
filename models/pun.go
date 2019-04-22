package models

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

// GeneratePuns given query word, homophone word list, and phrase
func GeneratePuns(word string, homophoneWords []WordRow, phrases []Phrase) []string {
	puns := []string{}

	for _, phrase := range phrases {
		tokens := strings.Split(phrase.PhraseText, " ")
		logrus.Infoln("Tokens", tokens)
		text := []string{}

		for _, token := range tokens {
			for _, homophoneWord := range homophoneWords {
				if strings.ToLower(token) == homophoneWord.Word {
					text = append(text, word)
					logrus.Infoln("Text Match", text)
				} else {
					text = append(text, token)
					logrus.Infoln("Text No Match", text)
				}
			}
		}
		result := strings.Join(text, " ")
		logrus.Infoln("Result", result)
		puns = append(puns, result)
		logrus.Infoln("Puns", puns)
	}

	return puns
}
