package main

import "strings"

func sanitize(s string) (string, bool) {
	s = strings.TrimSpace(s)

	// comment line or empty line
	if strings.HasPrefix(s, "//") || len(s) == 0 {
		return "", false
	}

	// in-line comments
	if strings.Contains(s, "//") {
		idx := strings.Index(s, "//")
		s = strings.TrimSpace(s[:idx])
	}

	return s, true
}
