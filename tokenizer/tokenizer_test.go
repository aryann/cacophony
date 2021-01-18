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
		{input: "(",
			want: []Token{{Type: LeftParen, Value: "("}}},
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
		/*
			{input: ")",
				want: []Token{{Type: RightParen}}},
			{input: "(()",
				want: []Token{{Type: LeftParen}, {Type: LeftParen}, {Type: RightParen}}},
			{input: "(f)",
				want: []Token{{Type: LeftParen}, {Type: Identifier, Value: "f"}, {Type: RightParen}}},
			{input: "(define (f x y z) (+ x y z))",
				want: []Token{{Type: LeftParen}, {Type: Identifier, Value: "define"},
					{Type: LeftParen}, {Type: Identifier, Value: "f"}, {Type: Identifier, Value: "x"},
					{Type: Identifier, Value: "y"}, {Type: Identifier, Value: "z"}, {Type: RightParen},
					{Type: LeftParen}, {Type: Identifier, Value: "+"}, {Type: Identifier, Value: "x"},
					{Type: Identifier, Value: "y"}, {Type: Identifier, Value: "z"}, {Type: RightParen},
					{Type: RightParen}}},
			{input: "(f)(g)",
				want: []Token{{Type: LeftParen}, {Type: Identifier, Value: "f"}, {Type: RightParen},
					{Type: LeftParen}, {Type: Identifier, Value: "g"}, {Type: RightParen}}},
			{input: `("hello")`,
				want: []Token{{Type: LeftParen}, {Type: String, Value: "hello"}, {Type: RightParen}}},
			{input: `"hello"`,
				want: []Token{{Type: String, Value: "hello"}}},
			{input: `"\"hello\""`,
				want: []Token{{Type: String, Value: `"hello"`}}},
			{input: `"\n"`,
				want: []Token{{Type: String, Value: "\n"}}},
			{input: `"\x"`,
				want: []Token{{Type: String, Value: `\x`}}},
			{input: `""`,
				want: []Token{{Type: String, Value: ""}}},
			{input: `"\""`,
				want: []Token{{Type: String, Value: `"`}}},
			{input: `(define (f x y z) ("hello"))`,
				want: []Token{{Type: LeftParen}, {Type: Identifier, Value: "define"},
					{Type: LeftParen}, {Type: Identifier, Value: "f"}, {Type: Identifier, Value: "x"},
					{Type: Identifier, Value: "y"}, {Type: Identifier, Value: "z"}, {Type: RightParen},
					{Type: LeftParen}, {Type: String, Value: "hello"}, {Type: RightParen},
					{Type: RightParen}}},*/
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
