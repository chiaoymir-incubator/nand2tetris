package main

import (
	"regexp"
	"strings"
)

type Sanitizer struct {
	s string
}

func (st *Sanitizer) Sanitize() (string, bool) {
	var ignore bool
	st.s, ignore = st.removeComments()
	return st.s, ignore
}

func (st *Sanitizer) removeComments() (string, bool) {
	st.s = strings.TrimSpace(st.s)

	// comment line or empty line
	if strings.HasPrefix(st.s, "//") || len(st.s) == 0 {
		return "", true
	}

	// in-line comments
	if strings.Contains(st.s, "//") {
		idx := strings.Index(st.s, "//")
		st.s = strings.TrimSpace(st.s[:idx])
	}

	// multi-line comments
	// https://stackoverflow.com/questions/12682405/strip-out-c-style-comments-from-a-byte
	re := regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")
	st.s = re.ReplaceAllString(st.s, "")

	return st.s, false
}
