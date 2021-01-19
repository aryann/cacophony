package tokenizer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

var key = map[string]Type{
	":true":   Boolean,
	":false":  Boolean,
	":define": Define,
	":if":     If,
}

const eof = -1

type Type int

const (
	EOF = iota
	Error
	LeftParen
	RightParen
	Number
	Identifier
	String
	Boolean
	Define
	If
)

func (t Type) String() string {
	return []string{
		"EOF",
		"Error",
		"LeftParen",
		"RightParen",
		"Number",
		"Identifier",
		"String",
		"Boolean",
		"Define",
		"If",
	}[t]
}

type Token struct {
	Type  Type
	Value string
	line  int
	col   int
}

func (t Token) String() string {
	return fmt.Sprintf("%s<%s>", t.Type, t.Value)
}

type tokenizer struct {
	buf        string
	start      int
	pos        int
	width      int
	parenDepth int
	tokens     []Token
}

func (t *tokenizer) emit(tokenType Type) {
	t.tokens = append(t.tokens, Token{
		Type:  tokenType,
		Value: t.buf[t.start:t.pos],
	})
	t.start = t.pos
}

func (t *tokenizer) ignore() {
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

func (t *tokenizer) peek() rune {
	r := t.next()
	t.backup()
	return r
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
	r := t.next()
	for unicode.IsSpace(r) {
		r = t.next()
	}
	if r == eof {
		return nil
	}
	t.backup()
	t.ignore()
	return lexBody
}

func lexBody(t *tokenizer) stateFn {
	r := t.next()
	switch {
	case r == '(':
		t.emit(LeftParen)
		t.parenDepth++
		return lexBody
	case r == ')':
		t.parenDepth--
		if t.parenDepth < 0 {
			return t.errorf("unexpected right paren")
		}
		t.emit(RightParen)
		return lexBody
	case r == '"':
		return lexString
	case isAlphaNumeric(r):
		t.backup()
		return lexIdentifier
	case r == ':':
		return lexBuiltIn
	case r == eof:
		if t.parenDepth != 0 {
			return t.errorf("unterminated left paren")
		}
		return nil
	case unicode.IsSpace(r):
		t.backup()
		return lexSpace
	default:
		return t.errorf("unexpected character: %v", r)
	}
}

func lexIdentifier(t *tokenizer) stateFn {
	for {
		r := t.next()
		if !isAlphaNumeric(r) {
			if r != eof {
				t.backup()
			}
			t.emit(Identifier)
			break
		}
	}
	return lexBody
}

func lexBuiltIn(t *tokenizer) stateFn {
	for {
		r := t.next()
		if !isAlphaNumeric(r) {
			if r != eof {
				t.backup()
			}
			word := t.buf[t.start:t.pos]
			keywordType, ok := key[word]
			if !ok {
				return t.errorf("unexpected keyword: %s", word[1:])
			}
			t.emit(keywordType)
			break
		}
	}
	return lexBody
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

func Tokenize(buf string) []Token {
	t := tokenizer{
		buf:    buf,
		tokens: make([]Token, 0),
	}
	for state := lexBody; state != nil; {
		state = state(&t)
	}
	return t.tokens
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
