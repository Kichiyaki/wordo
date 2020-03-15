package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/sqweek/dialog"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	reg, err := regexp.Compile("[^a-zA-ZąęśżźćńłóĄĘŚŻŹĆŃŁÓ]+")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loading file...")
	filename, err := dialog.File().Filter("PDF file", "pdf").Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Reading file...")
	output, err := readPdf(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Processing file content...")
	wordFrequencies := make(map[string]int)
	output = strings.Trim(output, "")
	for _, word := range strings.Split(output, " ") {
		processedString := strings.ToLower(reg.ReplaceAllString(word, ""))
		if strings.Trim(processedString, "") == "" {
			continue
		}
		wordFrequencies[processedString]++
	}
	response := ""
	for i, pair := range rankByWordCount(wordFrequencies) {
		if i >= 99 {
			break
		}
		response += fmt.Sprintf("%s:%d", pair.Key, pair.Value)
		if i != 98 {
			response += "\n"
		}
	}
	filename, err = dialog.File().Filter("txt file", "txt").Title("Export").Save()
	err = ioutil.WriteFile(filename, []byte(response), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
