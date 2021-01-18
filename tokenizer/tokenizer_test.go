package tokenizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input string
		want  []Token
	}{
		{input: "",
			want: []Token{}},
		{input: "     ",
			want: []Token{}},
		{input: "   \n  \n",
			want: []Token{}},

		{input: "(",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: Error, Value: "unterminated left paren"},
			}},
		{input: ")",
			want: []Token{{Type: Error, Value: "unexpected right paren"}}},
		{input: "()",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"}}},
		{input: "(())",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"},
				{Type: RightParen, Value: ")"},
			}},
		{input: "(() ())",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"},
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"},
				{Type: RightParen, Value: ")"},
			}},
		{input: "(   )",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"}}},
		{input: "   (   )     ",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: RightParen, Value: ")"}}},

		{input: `"`,
			want: []Token{{Type: Error, Value: "unterminated string"}}},
		{input: `""`,
			want: []Token{{Type: String, Value: `""`}}},
		{input: `"\t"`,
			want: []Token{{Type: String, Value: `"\t"`}}},
		{input: "\"\n\"",
			want: []Token{{Type: Error, Value: "unterminated string"}}},
		{input: `"\""`,
			want: []Token{{Type: String, Value: `"\""`}}},
		{input: `"hello world"`,
			want: []Token{{Type: String, Value: `"hello world"`}}},
		{input: `"hello \"world\""`,
			want: []Token{{Type: String, Value: `"hello \"world\""`}}},

		{input: "i",
			want: []Token{{Type: Identifier, Value: "i"}}},
		{input: "    i   ",
			want: []Token{{Type: Identifier, Value: "i"}}},
		{input: "identifier",
			want: []Token{{Type: Identifier, Value: "identifier"}}},
		{input: "   identifier",
			want: []Token{{Type: Identifier, Value: "identifier"}}},
		{input: "identifier    ",
			want: []Token{{Type: Identifier, Value: "identifier"}}},
		{input: "_identifier",
			want: []Token{{Type: Identifier, Value: "_identifier"}}},
		{input: "identifier651",
			want: []Token{{Type: Identifier, Value: "identifier651"}}},

		{input: "( identifier  )",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: Identifier, Value: "identifier"},
				{Type: RightParen, Value: ")"},
			}},

		{input: "( identifier (x y z) )",
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: Identifier, Value: "identifier"},
				{Type: LeftParen, Value: "("},
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
				{Type: Identifier, Value: "z"},
				{Type: RightParen, Value: ")"},
				{Type: RightParen, Value: ")"},
			}},

		{input: `( identifier "hello" "world" )`,
			want: []Token{
				{Type: LeftParen, Value: "("},
				{Type: Identifier, Value: "identifier"},
				{Type: String, Value: `"hello"`},
				{Type: String, Value: `"world"`},
				{Type: RightParen, Value: ")"},
			}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			tokens := Tokenize2(testCase.input)
			if !reflect.DeepEqual(tokens, testCase.want) {
				t.Fatalf("want tokens %+v, got %+v", testCase.want, tokens)
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr string
	}{
		{input: `"`, wantErr: "1:1: syntax error: unexpected end of file"},
		{input: `(define
"`, wantErr: "2:2: syntax error: unexpected end of file"},
		{input: `(")`, wantErr: "1:3: syntax error: unexpected end of file"},
		{input: "A", wantErr: "1:1: syntax error: unexpected character"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			tokens, err := Tokenize(strings.NewReader(testCase.input))
			if err == nil {
				t.Fatalf("want error '%s', got none", testCase.wantErr)
			}
			if err.Error() != testCase.wantErr {
				t.Fatalf("want error '%s', got '%s'", testCase.wantErr, err.Error())
			}
			if tokens != nil {
				t.Fatalf("want nil tokens, got %v", tokens)
			}

		})
	}
}
