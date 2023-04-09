package wasm

import (
	"bufio"
	"io"
	"log"
	"strings"
)

const Whitespace = " \t\r\n"

// Skip a multiline comment (; ;)
func SkipComment(text string) string {
	if strings.HasPrefix(text, "(;") {
		p := strings.Index(text, ";)")
		if p == -1 {
			panic("Unclosed (; ;) comment")
		}
		text = strings.TrimLeft(text[p+2:], Whitespace)
	}
	return text
}

// Reads non-whitespace token
func ReadToken(text string) (string, string) {
	text = SkipComment(text)

	token := ""
	r := bufio.NewReader(strings.NewReader(text))
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			break
		}

		token = token + string(ch)
	}
	return token, strings.TrimLeft(text[len(token):], Whitespace)
}

// Reads a string enclosed with ""
func ReadString(text string) (string, string) {
	text = SkipComment(text)

	token := ""
	r := bufio.NewReader(strings.NewReader(text))
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		token = token + string(ch)

		if ch == '"' && len(token) > 1 {
			break
		}
	}
	return token, strings.TrimLeft(text[len(token):], Whitespace)
}

// This reads an element enclosed with parenthesis.
// It also keeps track of speechmarks
func ReadElement(text string) (string, string) {
	text = SkipComment(text)

	bracketCount := 0
	inString := false
	//	el := ""

	current := 0

	r := bufio.NewReader(strings.NewReader(text))
	for {
		ch, _, err := r.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		current++

		if ch == '"' {
			inString = !inString
		}

		// Only care about bracks not inside a string.
		if !inString {
			if ch == '(' {
				bracketCount++
			}
			if ch == ')' {
				bracketCount--
			}
		}

		if bracketCount == 0 {
			break
		}
	}

	return text[:current], strings.TrimLeft(text[current:], Whitespace)
}
