package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	headingText := doc.Find("h1").First().Text()
	if headingText == "" {
		headingText = doc.Find("h2").First().Text()
	}
	return strings.Join(strings.Fields(headingText), " ")
}
