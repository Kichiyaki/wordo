package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/sqweek/dialog"
)

type Config struct {
	Top               int    `json:"top"`
	MinimumWordLength int    `json:"minimum_word_length"`
	Regex             string `json:"regex"`
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var cfg = &Config{}

func main() {
	log.Print("Loading config...")
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	log.Print("Parsing config...")
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteValue, cfg)
	if err != nil {
		log.Fatal(err)
	}
	reg, err := regexp.Compile(cfg.Regex)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Loading file...")
	filename, err := dialog.File().Filter("PDF file", "pdf").Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Reading file...")
	fileContent, err := readPdf(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Processing file content...")
	wordFrequencies := make(map[string]int)
	fileContent = strings.Trim(fileContent, "")
	total := 0
	for _, word := range strings.Split(fileContent, " ") {
		processedString := strings.Trim(strings.ToLower(reg.ReplaceAllString(word, "")), "")
		if processedString != "" && ((cfg.MinimumWordLength > 0 && len(processedString) >= cfg.MinimumWordLength) || cfg.MinimumWordLength <= 0) {
			wordFrequencies[processedString]++
			total++
		}
	}
	log.Print("Preparing output...")
	output := ""
	for i, pair := range rankByWordCount(wordFrequencies) {
		if i > cfg.Top-1 {
			break
		}
		output += fmt.Sprintf("%s;%d;%f", pair.Key, pair.Value, float64(pair.Value)/float64(total))
		if i != cfg.Top-1 {
			output += "\n"
		}
	}
	log.Print("Saving output...")
	filename, err = dialog.File().Filter("txt file", "txt").Title("Export").Save()
	err = ioutil.WriteFile(filename, []byte(output), 0644)
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
