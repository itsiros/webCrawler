package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func ExtractPageData(html, pageURL string) PageData {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}
	}

	urls, err := GetURLsFromHTML(html, baseURL)
	if err != nil {
		return PageData{}
	}

	imgs, err := GetImagesFromHTML(html, baseURL)
	if err != nil {
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		H1:             GetH1FromHTML(html),
		FirstParagraph: GetFirstParagraphFromHTML(html),
		OutgoingLinks:  urls,
		ImageURLs:      imgs,
	}
}

func GetH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1").First()
	if h1.Length() == 0 {
		return ""
	}

	return strings.TrimSpace(h1.Text())
}

func GetFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main").First()
	if main.Length() > 0 {
		p := main.Find("p").First()
		if p.Length() > 0 {
			return strings.TrimSpace(p.Text())
		}
	}

	p := doc.Find("p").First()
	if p.Length() == 0 {
		return ""
	}

	return strings.TrimSpace(p.Text())
}

func GetURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each((func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		abs := baseURL.ResolveReference(u)
		urls = append(urls, abs.String())

	}))

	if urls == nil {
		urls = []string{}
	}
	return urls, nil
}

func GetImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var imgs []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			return
		}

		abs := baseURL.ResolveReference(u)
		imgs = append(imgs, abs.String())
	})

	if imgs == nil {
		imgs = []string{}
	}
	return imgs, nil
}
