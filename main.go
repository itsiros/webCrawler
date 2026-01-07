package main

import (
	"fmt"
	"net/http"
	"os"
)

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("no website provided")
	}
	if len(os.Args) > 2 {
		return fmt.Errorf("too many arguments provided")
	}

	fmt.Println("starting crawl of:", os.Args[1])
	return nil
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
