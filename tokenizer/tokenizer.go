package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type Type int

const (
	EOF = iota
	Error
	LeftParen
	RightParen
	Number
	Identifier
	BuiltIn
	String
)

func (t Type) String() string {
	return []string{"EOF", "Error", "LeftParen", "RightParen", "Number", "Identifier", "BuiltIn", "String"}[t]
}

type Token struct {
	Type  Type
	Value string
	line  int
	col   int
}

type tokenizer struct {
	buf    string
	start  int
	pos    int
	width  int
	tokens []Token
}

func (t *tokenizer) emit(tokenType Type) {
	t.tokens = append(t.tokens, Token{
		Type:  tokenType,
		Value: t.buf[t.start:t.pos],
	})
	t.start = t.pos
}

func (t *tokenizer) next() rune {
	if t.pos >= len(t.buf) {
		return eof
	}
	result, width := utf8.DecodeRuneInString(t.buf[t.pos:])
	t.pos += width
	t.width = width
	return result
}

func (t *tokenizer) backup() {
	t.pos -= t.width
	t.width = 0
}

func (t *tokenizer) errorf(format string, args ...interface{}) stateFn {
	t.tokens = append(t.tokens, Token{
		Type:  Error,
		Value: fmt.Sprintf(format, args...),
	})
	return nil
}

type stateFn func(*tokenizer) stateFn

func lexSpace(t *tokenizer) stateFn {
	for unicode.IsSpace(t.next()) {
	}
	t.backup()
	return lexBody
}

func lexBody(t *tokenizer) stateFn {
	r := t.next()
	switch {
	case r == '(':
		t.emit(LeftParen)
		lexSpace(t)
		return lexFunction
	case r == '"':
		return lexString
	case isAlphaNumeric(r):
		t.backup()
		return lexIdentifier
	case r == eof:
		return nil
	default:
		return t.errorf("unexpected character: %v", r)
	}
}

func lexFunction(t *tokenizer) stateFn {
	return nil
}

func lexIdentifier(t *tokenizer) stateFn {
	return nil

}

func lexString(t *tokenizer) stateFn {
	for {
		switch r := t.next(); r {
		case '\\':
			if r := t.next(); r == eof || r == '\n' {
				return t.errorf("unterminated string")
			}
		case eof, '\n':
			return t.errorf("unterminated string")
		case '"':
			t.emit(String)
			return lexBody
		}
	}
}

func Tokenize2(buf string) []Token {
	t := tokenizer{
		buf:    buf,
		tokens: make([]Token, 0),
	}
	for state := lexSpace; state != nil; {
		state = state(&t)
	}
	return t.tokens
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (t Token) Error(message string, a ...interface{}) error {
	return fmt.Errorf("invalid syntax on line %d, column %d: %s",
		t.line, t.col, fmt.Sprintf(message, a...))
}

func isBuiltInStart(r rune) bool {
	return r == ':'
}

func isIdentifierStart(r rune) bool {
	return 'a' <= r && r <= 'z' || r == '+' || r == '-'
}

func isIdentifierChar(r rune) bool {
	return isIdentifierStart(r) || '0' <= r && r <= '9'
}

type scanner struct {
	buf *bytes.Buffer

	prevLine int
	prevCol  int
	currLine int
	currCol  int
}

func newScanner(buf []byte) *scanner {
	return &scanner{buf: bytes.NewBuffer(buf), prevLine: 1, prevCol: 1}
}

func (s *scanner) Peek() (rune, error) {
	next, _, err := s.buf.ReadRune()
	if err != nil {
		return next, err
	}
	if err := s.buf.UnreadRune(); err != nil {
		return rune(0), err
	}
	return next, nil
}

func (s *scanner) Next() (rune, error) {
	next, _, err := s.buf.ReadRune()
	if err != nil {
		return next, err
	}

	s.prevCol = s.currCol
	if next == '\n' {
		s.prevLine = s.currLine
		s.currLine++
		s.currCol = 1
	} else {
		s.currCol++
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

func (s *scanner) Error(message string, a ...interface{}) error {
	return fmt.Errorf("invalid syntax on line %d, column %d: %s",
		s.prevLine, s.prevCol, fmt.Sprintf(message, a...))
}

func (s *scanner) NewToken(t Type, value string) Token {
	return Token{
		Type:  t,
		Value: value,
		line:  s.prevLine,
		col:   s.prevCol,
	}
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
			tokens = append(tokens, scanner.NewToken(LeftParen, ""))
		} else if next == ')' {
			tokens = append(tokens, scanner.NewToken(RightParen, ""))

		} else if isIdentifierStart(next) || isBuiltInStart(next) {
			var builder strings.Builder

			var t Type
			if isIdentifierStart(next) {
				t = Identifier
				builder.WriteRune(next)
			} else {
				t = BuiltIn
			}

			for {
				next, err = scanner.Peek()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}

				if isIdentifierChar(next) {
					builder.WriteRune(next)
					if _, err := scanner.Next(); err != nil {
						return nil, err
					}
				} else {
					break
				}
			}

			tokens = append(tokens, scanner.NewToken(t, builder.String()))

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
			tokens = append(tokens, scanner.NewToken(String, builder.String()))

		} else {
			return nil, scanner.Error("unexpected character: %c", next)
		}
	}

	return tokens, nil
}
