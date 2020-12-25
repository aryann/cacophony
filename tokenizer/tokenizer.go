package tokenizer

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Type int

const (
	LeftParen = iota
	RightParen
	Number
	Identifier
	String
)

func (t Type) String() string {
	return []string{"LeftParen", "RightParen", "Number", "Identifier", "String"}[t]
}

type Token struct {
	Type  Type
	Value string
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

		} else if rune == '"' {
			var builder strings.Builder
			limit := i + 1
			for limit <= len(buf) {
				rune, width = utf8.DecodeRune(buf[limit:])

				if rune == '\\' {
					if limit+1 >= len(buf) {
						return nil, fmt.Errorf("illegal")
					}
					limit += width
					rune, width = utf8.DecodeRune(buf[limit:])
					switch rune {
					case '"':
						builder.WriteRune('"')
					case 'n':
						builder.WriteRune('\n')
					case '\\':
						builder.WriteRune('\\')
					default:
						builder.WriteRune('\\')
						builder.WriteRune(rune)
					}

				} else if rune == '"' {
					break
				} else {
					builder.WriteRune(rune)
				}

				limit += width
			}
			tokens = append(tokens, Token{Type: String, Value: builder.String()})
			i = limit + 1

		} else {
			return nil, fmt.Errorf("illegal character at offset %d", i)
		}
	}

	return tokens, nil
}
