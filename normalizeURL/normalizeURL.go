package normalizeURL

import (
	"net/url"
	"strings"
)

func NormalizeURL(urlString string) string {

	u, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	var b strings.Builder

	if u.Host != "" {
		b.WriteString(u.Host)
	}

	if u.Path != "/" {
		b.WriteString(u.Path)
	}

	if u.RawQuery != "" {
		b.WriteByte('?')
		b.WriteString(u.RawQuery)
	}

	if u.Fragment != "" {
		b.WriteByte('#')
		b.WriteString(u.Fragment)
	}

	return b.String()
}
