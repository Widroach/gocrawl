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
	headingText := doc.Find("h1, h2").First().Text()
	return trim(headingText)
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	firstParagraph := doc.Find("main").First().Find("p").First()
	if trim(firstParagraph.Text()) == "" {
		firstParagraph = doc.Find("p").First()
	}

	return trim(firstParagraph.Text())
}

func trim(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
