package getfromhtml

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single h1",
			input:    "<html><body><h1>Test Title</h1></body></html>",
			expected: "Test Title",
		},
		{
			name:     "multiple h1 uses first",
			input:    "<html><body><h1>First</h1><h1>Second</h1></body></html>",
			expected: "First",
		},
		{
			name:     "no h1",
			input:    "<html><body><p>No header</p></body></html>",
			expected: "",
		},
		{
			name:     "nested tags inside h1",
			input:    "<html><body><h1>Hello <span>World</span></h1></body></html>",
			expected: "Hello World",
		},
		{
			name:     "whitespace trimmed",
			input:    "<html><body><h1>   Trim Me   </h1></body></html>",
			expected: "Trim Me",
		},
		{
			name:     "uppercase h1 tag",
			input:    "<html><body><H1>Uppercase</H1></body></html>",
			expected: "Uppercase",
		},
		{
			name:     "h1 outside body",
			input:    "<h1>Outside</h1>",
			expected: "Outside",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getH1FromHTML(tc.input)

			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "main paragraph has priority",
			input: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "multiple paragraphs inside main uses first",
			input: `<html><body>
				<main>
					<p>First main paragraph.</p>
					<p>Second main paragraph.</p>
				</main>
			</body></html>`,
			expected: "First main paragraph.",
		},
		{
			name: "fallback to body paragraph when no main",
			input: `<html><body>
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
			</body></html>`,
			expected: "First paragraph.",
		},
		{
			name: "main exists but has no paragraph fallback to body",
			input: `<html><body>
				<main>
					<div>No paragraph here</div>
				</main>
				<p>Body paragraph.</p>
			</body></html>`,
			expected: "Body paragraph.",
		},
		{
			name: "paragraph inside main with nested tags",
			input: `<html><body>
				<main>
					<p>Hello <strong>World</strong></p>
				</main>
			</body></html>`,
			expected: "Hello World",
		},
		{
			name:     "paragraph outside body",
			input:    `<p>Lonely paragraph.</p>`,
			expected: "Lonely paragraph.",
		},
		{
			name:     "no paragraph anywhere",
			input:    `<html><body><div>No content</div></body></html>`,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.input)

			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	inputURL := "https://blog.boot.dev"

	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "Absolute URL",
			html:     `<a href="https://blog.boot.dev/page1">Page1</a>`,
			expected: []string{"https://blog.boot.dev/page1"},
		},
		{
			name:     "Relative URL",
			html:     `<a href="/page2">Page2</a>`,
			expected: []string{"https://blog.boot.dev/page2"},
		},
		{
			name:     "Multiple links",
			html:     `<a href="/one">One</a><a href="https://other.com/two">Two</a>`,
			expected: []string{"https://blog.boot.dev/one", "https://other.com/two"},
		},
		{
			name:     "Missing href",
			html:     `<a>No link</a>`,
			expected: []string{},
		},
	}

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("couldn't parse input URL: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tt.html, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

/* -------------------- IMAGE TESTS -------------------- */

func TestGetImagesFromHTML(t *testing.T) {
	inputURL := "https://blog.boot.dev"

	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "Relative image",
			html:     `<img src="/logo.png" alt="Logo">`,
			expected: []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name:     "Absolute image",
			html:     `<img src="https://example.com/image.jpg">`,
			expected: []string{"https://example.com/image.jpg"},
		},
		{
			name:     "Multiple images",
			html:     `<img src="/img1.png"><img src="https://cdn.com/img2.jpg">`,
			expected: []string{"https://blog.boot.dev/img1.png", "https://cdn.com/img2.jpg"},
		},
		{
			name:     "Missing src",
			html:     `<img alt="No src">`,
			expected: []string{},
		},
	}

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("couldn't parse input URL: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getImagesFromHTML(tt.html, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	// ---------- Case 1: Basic page ----------
	{
		inputURL := "https://blog.boot.dev"
		inputBody := `<html><body>
			<h1>Test Title</h1>
			<p>This is the first paragraph.</p>
			<a href="/link1">Link 1</a>
			<img src="/image1.jpg" alt="Image 1">
		</body></html>`

		actual := extractPageData(inputBody, inputURL)

		expected := PageData{
			URL:            "https://blog.boot.dev",
			H1:             "Test Title",
			FirstParagraph: "This is the first paragraph.",
			OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
			ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("basic case failed: expected %+v, got %+v", expected, actual)
		}
	}

	// ---------- Case 2: No links or images ----------
	{
		inputURL := "https://blog.boot.dev"
		inputBody := `<html><body>
			<h1>Only Title</h1>
			<p>Just text.</p>
		</body></html>`

		actual := extractPageData(inputBody, inputURL)

		expected := PageData{
			URL:            "https://blog.boot.dev",
			H1:             "Only Title",
			FirstParagraph: "Just text.",
			OutgoingLinks:  []string{},
			ImageURLs:      []string{},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("no links/images case failed: expected %+v, got %+v", expected, actual)
		}
	}

	// ---------- Case 3: Missing H1 ----------
	{
		inputURL := "https://blog.boot.dev"
		inputBody := `<html><body>
			<p>Paragraph without title.</p>
			<a href="/link">Link</a>
		</body></html>`

		actual := extractPageData(inputBody, inputURL)

		expected := PageData{
			URL:            "https://blog.boot.dev",
			H1:             "",
			FirstParagraph: "Paragraph without title.",
			OutgoingLinks:  []string{"https://blog.boot.dev/link"},
			ImageURLs:      []string{},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("missing h1 case failed: expected %+v, got %+v", expected, actual)
		}
	}

	// ---------- Case 4: Multiple links and images ----------
	{
		inputURL := "https://blog.boot.dev"
		inputBody := `<html><body>
			<h1>Main</h1>
			<p>First paragraph.</p>
			<a href="/one">One</a>
			<a href="https://example.com/two">Two</a>
			<img src="/img1.png">
			<img src="https://cdn.com/img2.png">
		</body></html>`

		actual := extractPageData(inputBody, inputURL)

		expected := PageData{
			URL:            "https://blog.boot.dev",
			H1:             "Main",
			FirstParagraph: "First paragraph.",
			OutgoingLinks: []string{
				"https://blog.boot.dev/one",
				"https://example.com/two",
			},
			ImageURLs: []string{
				"https://blog.boot.dev/img1.png",
				"https://cdn.com/img2.png",
			},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("multiple links/images case failed: expected %+v, got %+v", expected, actual)
		}
	}

	// ---------- Case 5: Empty HTML ----------
	{
		inputURL := "https://blog.boot.dev"
		inputBody := ``

		actual := extractPageData(inputBody, inputURL)

		expected := PageData{
			URL:            "https://blog.boot.dev",
			H1:             "",
			FirstParagraph: "",
			OutgoingLinks:  []string{},
			ImageURLs:      []string{},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("empty html case failed: expected %+v, got %+v", expected, actual)
		}
	}
}
