package main

import (
	"strconv"
	"strings"
)

type TokenType string

const (
	TokenInvalid     TokenType = "invalid"
	TokenKeyword     TokenType = "keyword"
	TokenSymbol      TokenType = "symbol"
	TokenIntConst    TokenType = "integerConstant"
	TokenStringConst TokenType = "stringConstant"
	TokenIdentifier  TokenType = "identifier"
)

var (
	tc = TokenChecker{}
)

type Token struct {
	isValid   bool
	tokenType TokenType
	content   string
}

type TokenChecker struct {
}

func (tc *TokenChecker) checkToken(s string) (TokenType, bool) {
	switch {
	case len(s) == 0:
		return TokenInvalid, false

	case tc.isKeyword(s):
		return TokenKeyword, true

	case tc.isSymbol(s[0]):
		if len(s) != 1 {
			return TokenInvalid, false
		}
		return TokenSymbol, true

	case strings.Contains(s, "\""):
		if strings.Count(s, "\"") != 2 || !strings.HasPrefix(s, "\"") || !strings.HasSuffix(s, "\"") {
			return TokenInvalid, false
		}
		return TokenStringConst, true

	default:
		if _, err := strconv.Atoi(string(s[0])); err != nil {
			// start with non-numeric character
			return TokenIdentifier, true
		} else if i, err := strconv.Atoi(s); err == nil && i >= 0 && i <= 32767 {
			// the string is an int constant
			return TokenIntConst, true
		}
		return TokenInvalid, false
	}
}

func (tc *TokenChecker) isIgnoredChar(b byte) bool {
	switch b {
	case ' ', '\t', '\n':
		return true
	}
	return false
}

func (tc *TokenChecker) isKeyword(s string) bool {
	switch s {
	case "class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return":
		return true
	}
	return false
}

func (tc *TokenChecker) isSymbol(b byte) bool {
	switch b {
	case '(', ')', '[', ']', '{', '}', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~':
		return true
	}
	return false
}

type Tokenizer struct {
	str   string
	idx   int
	token Token
}

func NewTokenizer(s string) Tokenizer {
	tokenizer := Tokenizer{}
	tokenizer.str = s
	tokenizer.token = Token{}
	return tokenizer
}

func (tknzr *Tokenizer) isEnd() bool {
	return tknzr.idx >= len(tknzr.str)
}

func (tknzr *Tokenizer) Advance() {
	// mark the tokenizer as ending
	if tknzr.isEnd() {
		tknzr.token = Token{}
		return
	}

	var sb strings.Builder
	var quoteCount int

	// append char until it meet ignore character
	// if we are in a quote, we have to intentionally break the loop until we neet the other quote
	for !tknzr.isEnd() && (quoteCount > 0 || !tc.isIgnoredChar(tknzr.str[tknzr.idx])) {
		c := tknzr.str[tknzr.idx]
		// if the builder is not empty and the current char is a symbol
		// we can know there is an identifier before to be parsed
		if sb.Len() != 0 && tc.isSymbol(c) && quoteCount == 0 {
			break
		}

		sb.WriteByte(c)
		tknzr.idx++

		if c == '"' {
			quoteCount++
		}

		// we meet the other quote
		if quoteCount >= 2 {
			break
		}

		// detect a keyword
		if tc.isKeyword(sb.String()) && (tknzr.idx+1 < len(tknzr.str)) && tknzr.str[tknzr.idx+1] == ' ' {
			break
		}

		// detect a symbol
		if sb.Len() == 1 && tc.isSymbol(sb.String()[0]) {
			break
		}
	}

	// parse token
	if sb.Len() != 0 {
		if tokenType, ok := tc.checkToken(sb.String()); ok {
			tknzr.token = Token{
				isValid:   true,
				tokenType: tokenType,
			}

			if tokenType == TokenStringConst {
				s, _ := strconv.Unquote(sb.String())
				tknzr.token.content = s
			} else if tokenType == TokenSymbol {
				tknzr.token.content = sb.String()
				/* for output to xml */
				// switch sb.String() {
				// case "<":
				// 	tkzn.token.content = "&lt;"
				// case ">":
				// 	tkzn.token.content = "&gt;"
				// case "\"":
				// 	tkzn.token.content = "&quot;"
				// case "&":
				// 	tkzn.token.content = "&amp;"
				// default:
				// 	tkzn.token.content = sb.String()
				// }
			} else {
				tknzr.token.content = sb.String()
			}

		}
	}

	// remove ignored characters
	for !tknzr.isEnd() && tc.isIgnoredChar(tknzr.str[tknzr.idx]) {
		tknzr.idx++
	}
}

func (tknzr *Tokenizer) HasMoreTokens() bool {
	return tknzr.token.isValid
}

func (tknzr *Tokenizer) GetToken() Token {
	return tknzr.token
}
