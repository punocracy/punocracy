package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Parse data for loading into MySQL")

	// Load CSV data into memory
	csvFile, err := os.Open("../data/homophones.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Create a scanner from csvReader, and split by line
	csvScanner := bufio.NewScanner(csvFile)
	csvScanner.Split(bufio.ScanLines)

	// Scan CSV data into rawData line by line
	var rawData [][]string
	for csvScanner.Scan() {
		// Get one line from the scanner and split by commas
		lineText := csvScanner.Text()
		lineData := strings.Split(lineText, ",")

		// Append to data matrix
		rawData = append(rawData, lineData)
	}

	// Check for scan errors
	if err := csvScanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create map of words to group numbers. For each line (homophone group), insert
	// 	all memebers of the group with the same group number. When a repeat of the
	// 	line is found, words are not re-inserted since they already exist in the map
	var groupNum int                // Homophone group identifier (integer)
	wordMap := make(map[string]int) // Map of words to group number
	for _, group := range rawData {
		// If the group does not exist in the map, then insert all its members with
		// 	the same groupNum; increment groupNum
		if _, groupExists := wordMap[group[0]]; !groupExists {
			// Insert all members into the map with the same groupNum
			for _, word := range group {
				wordMap[word] = groupNum
			}
			groupNum = groupNum + 1
		}
	}

	// Iterate through map to find all words in groups 0 - 3
	// TODO: remove this debugging code
	for key, val := range wordMap {
		if val >= 0 && val <= 3 {
			fmt.Println(val, key)
		}
	}

	// Open output file
	outputCsv, err := os.Create("../data/parsedHomophones.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputCsv.Close()

	// Write line by line
	for word, groupNum := range wordMap {
		// Create output CSV line
		writeStr := fmt.Sprintf("%s,%d\n", word, groupNum)
		outputCsv.WriteString(writeStr)
	}
}
