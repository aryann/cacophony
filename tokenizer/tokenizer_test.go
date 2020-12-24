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
		{input: "(",
			want: []Token{{Type: LeftParen}}},
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			tokens, err := Tokenize(strings.NewReader(testCase.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tokens, testCase.want) {
				t.Fatalf("want tokens %+v, got %+v", testCase.want, tokens)
			}
		})
	}
}
