package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func run() error {

	if len(os.Args) != 4 {
		return fmt.Errorf("The crawler is used as follows: go run . \"domain.com\" MaxThreads MaxPages")
	}

	fmt.Println("starting crawl of:", os.Args[1])
	return nil
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Crawler/1.0")

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: status %d", resp.StatusCode)
	}

	contType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contType, "text/html") {
		return "", errors.New("content type is not html")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if maxConcurrency <= 0 || maxPages <= 0 {
		fmt.Fprintln(os.Stderr, "MaxThreads and MaxPages must be positive integers")
		return
	}

	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	if err = writeCSVReport(cfg.pages, "report.csv"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("CSV report generated")
}
