package main

import (
	"net/url"
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

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	urls := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return urls, err
	}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			if !isAbsolute(href) {
				href = resolveURL(*baseURL, href)
			}
			urls = append(urls, href)
		}
	})
	return urls, nil
}

func resolveURL(baseURL url.URL, relativePath string) (absoluteURL string) {
	relativeURL, _ := url.Parse(relativePath)
	return baseURL.ResolveReference(relativeURL).String()
}

func isAbsolute(href string) bool {
	u, err := url.Parse(href)
	if err != nil {
		return false
	}
	return u.IsAbs()
}

func trim(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
