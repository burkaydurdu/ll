package cmd

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"github.com/PuerkitoBio/goquery"
)

type Phrase struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Category string `json:"category"`
	Type     string `json:"type"`
}

type Phrases []Phrase

var phrases Phrases

func FetchFromTureng(query string) (Phrases, error) {
	var path string

	if viper.GetString("mode") == "Turkish - English" {
		path = "https://tureng.com/tr/turkce-ingilizce/"
	} else {
		path = "https://tureng.com/en/turkish-english/"
	}

	res, err := http.Get(path + query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New("status code error: %s " + res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	source := doc.Find("#englishResultsTable").Find("tbody tr")

	source.Each(func(i int, str *goquery.Selection) {
		var phrase Phrase

		tableHeaders := str.Find("th")
		tableColumns := str.Find("td")

		if tableHeaders.Eq(2).Text() != "" {
			phrase.Category = tableHeaders.Eq(1).Text()
			phrase.Source = tableHeaders.Eq(2).Text()
			phrase.Target = tableHeaders.Eq(3).Text()
			phrase.Type = control(viper.GetString("mode") == "Turkish - English", "Tür", "Type")
			phrases = append(phrases, phrase)
		} else {

			secondColumnTypeText := tableColumns.Eq(2).Find("i").Text()

			if secondColumnTypeText != "" {
				phrase.Type = convertType(secondColumnTypeText)
			} else {
				phrase.Type = convertType(tableColumns.Eq(3).Find("i").Text())
			}

			phrase.Category = tableColumns.Eq(1).Text()
			phrase.Source = tableColumns.Eq(2).Find("a").Text()
			phrase.Target = tableColumns.Eq(3).Find("a").Text()

			if phrase.Target != "" || phrase.Source != "" {
				phrases = append(phrases, phrase)
			}
		}
	})

	defer func() {
		phrases = []Phrase{}
	}()

	return phrases, nil
}

func convertType(phraseType string) string {
	phraseType = strings.Trim(phraseType, " ")
	switch phraseType {
	case "i.":
		return "isim"
	case "n.":
		return "noun"
	case "f.":
		return "fiil"
	case "v.":
		return "verb"
	case "zf.":
		return "zarf"
	case "adv.":
		return "adverb"
	case "s.":
		return "sıfat"
	case "adj.":
		return "adjective"
	case "ünl.":
		return "ünlem"
	case "interj.":
		return "interjection"
	default:
		return "unknown"
	}
}
