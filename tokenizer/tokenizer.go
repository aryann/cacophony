package tokenizer

import (
	"fmt"
	"io"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
)

type Type int

const (
	LeftParen = iota
	RightParen
	Number
	Identifier
)

func (t Type) String() string {
	return []string{"LeftParen", "RightParen", "Number", "Identifier"}[t]
}

type Token struct {
	Type  Type
	Value string
}

func (t Token) String() string {
	switch t.Type {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case Number, Identifier:
		return t.Value
	}
	return ""
}

func isIdentifierStart(r rune) bool {
	return 'a' <= r && r <= 'z' || r == '+' || r == '-'
}

func isIdentifierChar(r rune) bool {
	return isIdentifierStart(r) || '0' <= r && r <= '9'
}

func Tokenize(reader io.Reader) ([]Token, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	tokens := make([]Token, 0)
	i := 0

	for i < len(buf) {
		rune, width := utf8.DecodeRune(buf[i:])

		if unicode.IsSpace(rune) {
			i += width
		} else if rune == '(' {
			tokens = append(tokens, Token{Type: LeftParen})
			i += width
		} else if rune == ')' {
			tokens = append(tokens, Token{Type: RightParen})
			i += width

		} else if isIdentifierStart(rune) {
			start := i
			limit := i
			for limit <= len(buf) {
				rune, width = utf8.DecodeRune(buf[limit:])
				if isIdentifierChar(rune) {
					limit += width
				} else {
					break
				}
			}
			tokens = append(tokens, Token{Type: Identifier, Value: string(buf[start:limit])})
			i = limit

		} else {
			return nil, fmt.Errorf("illegal character at offset %d", i)
		}
	}

	return tokens, nil
}
