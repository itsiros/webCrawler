package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	defer wr.Flush()

	header := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if err = wr.Write(header); err != nil {
		return err
	}
	
	for url, page := range pages {
		row := []string{
			url,
			page.H1,
			page.FirstParagraph,
			strings.Join(page.OutgoingLinks, ";"),
			strings.Join(page.ImageURLs, ";"),
		}

		if err = wr.Write(row); err != nil {
			return err
		}
	}

	return nil
}
