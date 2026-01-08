# webCrawler

A lightweight, concurrent web crawler written in Go. It starts from a given URL, recursively crawls web pages, extracts and normalizes links, and generates detailed CSV reports on crawled pages and links.

## Overview

webCrawler is designed for educational and small-scale web scraping tasks. It demonstrates core concepts in concurrent programming, HTTP handling, HTML parsing, and data export. Key features include:
- Concurrent crawling using goroutines for efficiency
- URL normalization to avoid duplicates
- Configurable depth limits and concurrency levels
- CSV report generation for crawled pages and extracted links
- Robust error handling and logging

Built as a learning project in Go networking and data processing.

## Features

- ✅ Concurrent crawling with configurable goroutines
- ✅ URL normalization (removes fragments, handles relative URLs)
- ✅ HTML link extraction from `<a>` tags
- ✅ Depth-limited crawling to prevent infinite loops
- ✅ CSV export for pages (URL, status, links found) and links (source, target)
- ✅ Command-line configuration (URL, max pages, concurrency)
- ✅ Unit tests for core functions (URL normalization, HTML parsing)
- ✅ Error handling for network timeouts and invalid responses

## Project Structure

```
webCrawler/
├── go.mod                    # Go module definition
├── README.md                 # This file
├── cmd/
│   ├── main.go               # Entry point and CLI handling
│   ├── configure.go          # Configuration parsing and validation
│   ├── csvReport.go          # CSV report generation
│   ├── getFromHTML.go        # HTML parsing and link extraction
│   ├── getFromHTML_test.go   # Unit tests for HTML parsing
│   ├── normalizeURL.go       # URL normalization logic
│   └── normalizeURL_test.go  # Unit tests for URL normalization
└── test/                     # (Optional: test fixtures or sample data)
```

## Requirements

- **Go 1.25.4** or later
- Linux/macOS/Windows (tested on Linux)
- No external dependencies (uses only standard library: `net/http`, `golang.org/x/net/html`, etc.)

### Install Go (if needed):
```bash
# On Arch Linux
sudo pacman -S go

# On Ubuntu/Debian
sudo apt-get install golang-go

# Or download from https://golang.org/dl/
```

## Installation & Build

### Clone and build:
```bash
git clone https://github.com/itsiros/webCrawler.git
cd webCrawler

# Build the executable
go build -o webCrawler ./cmd

# Or run directly
go run ./cmd
```

### Run tests:
```bash
# Run all unit tests
go test ./cmd

# Run with verbose output
go test -v ./cmd

# Run specific test file
go test -v ./cmd -run TestNormalizeURL
```

## Usage

### Basic Crawl
Start crawling from a URL and generate reports:

```bash
# Crawl a site with defaults (max 100 pages, 5 concurrent workers)
./webCrawler https://example.com 5 50

### Command-Line Options
```bash
./webCrawler  <starting_url> [options]

Options:
  -maxPages int     Maximum pages to crawl (default 100)
  -concurrency int  Number of concurrent workers (default 5)
```

## Configuration

Configuration is handled via command-line flags (see Usage). For advanced setups, modify `cmd/configure.go` to add environment variables or config files.

- **Starting URL**: Required positional argument
- **Max Pages**: Limits crawl size to prevent resource exhaustion
- **Concurrency**: Number of goroutines; increase for faster crawling but monitor system load

## Examples

### Crawl a Small Site
```bash
./webCrawler https://golang.org 5 10 
# Crawls up to 10 pages from golang.org, generates reports in root folder
```

### High-Concurrency Crawl
```bash
./webCrawler https://example.com -concurrency 20 -maxPages 500
# Uses 20 workers, crawls up to 500 pages
```

## Architecture

### Core Components

- **main.go**: CLI entry point, orchestrates crawling and reporting
- **configure.go**: Parses flags and validates inputs
- **getFromHTML.go**: Parses HTML and extracts `<a href>` links
- **normalizeURL.go**: Cleans URLs (removes fragments, resolves relatives)
- **csvReport.go**: Writes structured CSV data

### Design Patterns

- **Worker Pool**: Goroutines for concurrent page fetching
- **Channel-Based Communication**: For distributing work and collecting results
- **Error Propagation**: Graceful handling of network failures
- **Modular Parsing**: Separate functions for URL and HTML processing

## Performance

- **Concurrent**: Scales with goroutines 
- **Memory Efficient**: Streams HTML parsing, no full page storage
- **Configurable Limits**: Prevents runaway crawling

**Approximate Benchmarks** (on test machine):
- Page fetch + parse: ~50-200ms per page (depends on site)
- Memory usage: <50MB for 1000 pages
- Startup time: <100ms

## Limitations

- No JavaScript rendering (static HTML only)
- Basic rate limiting (add delays if needed)
- No authentication or session handling
- Limited to HTTP/HTTPS (no FTP, etc.)
- Educational scope (not production-grade scraper)

## Future Enhancements

- [ ] Add rate limiting and politeness delays
- [ ] Support for robots.txt parsing
- [ ] JSON output option
- [ ] Web UI for visualization
- [ ] Docker containerization

## License

MIT License - see LICENSE file for details.

**Version:** 1.0.0  
**Go Version:** 1.25.4  
**Last Updated:** January 8, 2026