package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating the request: %v", err)
	}
	req.Header.Add("User-Agent", "GoCrawl/1.0.0")
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting %v: %v", rawURL, err)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", errors.New("resource returned non-html content")
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to request %v, returned %v", rawURL, resp.Status)
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the response: %v", err)
	}

	return string(html), nil
}

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
				href, err = resolveURL(*baseURL, href)
				if err != nil {
					return
				}
			}
			urls = append(urls, href)
		}
	})
	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	imageURLs := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return imageURLs, err
	}
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("src")
		if ok {
			if !isAbsolute(href) {
				href, err = resolveURL(*baseURL, href)
				if err != nil {
					return
				}
			}
			imageURLs = append(imageURLs, href)
		}
	})
	return imageURLs, nil
}

func extractPageData(html, pageURL string) PageData {
	// TODO: It doesn't return errors, it should be handled
	parsedBasedURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}
	}

	heading := getHeadingFromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)

	links, err := getURLsFromHTML(html, parsedBasedURL)
	if err != nil {
		return PageData{}
	}

	imageLinks, err := getImagesFromHTML(html, parsedBasedURL)
	if err != nil {
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  links,
		ImageURLs:      imageLinks,
	}
}

// This function safely constructs an absolute URL from the base URL and from the relative path. Warning, it assumes the relative path is correct when parsed with [url.Parse]
func resolveURL(baseURL url.URL, relativePath string) (absoluteURL string, err error) {
	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(relativeURL).String(), nil
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
