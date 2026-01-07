package normalizeURL_test

import (
	"testing"

	"github.com/itsiros/webCrawler/normalizeURL"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove https scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove http scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "no scheme present",
			inputURL: "blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "root domain only",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev",
		},
		{
			name:     "trailing slash removed",
			inputURL: "https://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
		{
			name:     "nested path preserved",
			inputURL: "https://blog.boot.dev/path/to/page",
			expected: "blog.boot.dev/path/to/page",
		},
		{
			name:     "query string preserved",
			inputURL: "https://blog.boot.dev/path?q=go",
			expected: "blog.boot.dev/path?q=go",
		},
		{
			name:     "fragment preserved",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path#section",
		},
		{
			name:     "query and fragment preserved",
			inputURL: "https://blog.boot.dev/path?q=go#section",
			expected: "blog.boot.dev/path?q=go#section",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := normalizeURL.NormalizeURL(tc.inputURL)
			if actual != tc.expected {
				t.Errorf(
					"Test %v - %s FAIL: expected URL: %v, actual: %v",
					i,
					tc.name,
					tc.expected,
					actual,
				)
			}
		})
	}
}
