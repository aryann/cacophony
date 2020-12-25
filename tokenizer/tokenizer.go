package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
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

type scanner struct {
	buf      *bytes.Buffer
	currLine int
	currCol  int
}

func newScanner(buf []byte) *scanner {
	return &scanner{buf: bytes.NewBuffer(buf), currLine: 1, currCol: 0}
}

func (s *scanner) Next() (rune, error) {
	next, _, err := s.buf.ReadRune()
	if err != nil {
		return rune(0), err
	}

	if next == '\n' {
		s.currLine++
		s.currCol = 0
	}
	return next, nil
}

func (t *scanner) ExpectNext() (rune, error) {
	next, err := t.Next()
	if err == io.EOF {
		return rune(0), t.Error("unexpected end of file")
	}
	return next, err
}

func (s *scanner) Error(message string) error {
	return fmt.Errorf("%d:%d: syntax error: %s", s.currLine, s.currCol, message)
}

func (s *scanner) Prev() error {
	return s.buf.UnreadRune()
}

func Tokenize(reader io.Reader) ([]Token, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	scanner := newScanner(buf)
	tokens := make([]Token, 0)

	for {
		next, err := scanner.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if unicode.IsSpace(next) {
			// Nothing to do!
		} else if next == '(' {
			tokens = append(tokens, Token{Type: LeftParen})
		} else if next == ')' {
			tokens = append(tokens, Token{Type: RightParen})

		} else if isIdentifierStart(next) {
			var builder strings.Builder
			builder.WriteRune(next)
			for {
				next, err = scanner.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}

				if isIdentifierChar(next) {
					builder.WriteRune(next)
				} else {
					if err := scanner.Prev(); err != nil {
						return nil, err
					}
					break
				}
			}

			tokens = append(tokens, Token{Type: Identifier, Value: builder.String()})

		} else if next == '"' {
			var builder strings.Builder
			for {
				next, err = scanner.ExpectNext()
				if err != nil {
					return nil, err
				}

				if next == '\\' {
					next, err = scanner.ExpectNext()
					if err != nil {
						return nil, err
					}

					switch next {
					case '"':
						builder.WriteRune('"')
					case 'n':
						builder.WriteRune('\n')
					case '\\':
						builder.WriteRune('\\')
					default:
						builder.WriteRune('\\')
						builder.WriteRune(next)
					}

				} else if next == '"' {
					break
				} else {
					builder.WriteRune(next)
				}
			}
			tokens = append(tokens, Token{Type: String, Value: builder.String()})

		} else {
			return nil, scanner.Error("unexpected character")
		}
	}

	return tokens, nil
}
