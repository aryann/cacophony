package tokenizer

import (
	"reflect"
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

		{input: "   ( identifier (x y z) )  ",
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

		{input: ":true",
			want: []Token{{Type: Boolean, Value: ":true"}}},
		{input: ":false",
			want: []Token{{Type: Boolean, Value: ":false"}}},
		{input: ":define",
			want: []Token{{Type: Define, Value: ":define"}}},
		{input: ":if",
			want: []Token{{Type: If, Value: ":if"}}},
		{input: ":unknown",
			want: []Token{{Type: Error, Value: "unexpected keyword: unknown"}}},

		{input: "1",
			want: []Token{{Type: Number, Value: "1"}}},
		{input: "1  ",
			want: []Token{{Type: Number, Value: "1"}}},
		{input: "  1",
			want: []Token{{Type: Number, Value: "1"}}},
		{input: "0",
			want: []Token{{Type: Number, Value: "0"}}},
		{input: "0.1",
			want: []Token{{Type: Number, Value: "0.1"}}},
		{input: "1.0",
			want: []Token{{Type: Number, Value: "1.0"}}},
		{input: "-0.1",
			want: []Token{{Type: Number, Value: "-0.1"}}},
		{input: "+0.1",
			want: []Token{{Type: Number, Value: "+0.1"}}},
		{input: "123",
			want: []Token{{Type: Number, Value: "123"}}},
		{input: "+123",
			want: []Token{{Type: Number, Value: "+123"}}},
		{input: "-123",
			want: []Token{{Type: Number, Value: "-123"}}},
		{input: "12.33",
			want: []Token{{Type: Number, Value: "12.33"}}},
		{input: "+12.33",
			want: []Token{{Type: Number, Value: "+12.33"}}},
		{input: "-12.33",
			want: []Token{{Type: Number, Value: "-12.33"}}},
		{input: "0x1234567890abcdef",
			want: []Token{{Type: Number, Value: "0x1234567890abcdef"}}},
		{input: "0X1",
			want: []Token{{Type: Number, Value: "0X1"}}},
		{input: "+0x1234567890abcdef",
			want: []Token{{Type: Number, Value: "+0x1234567890abcdef"}}},
		{input: "-0x1234567890abcdef",
			want: []Token{{Type: Number, Value: "-0x1234567890abcdef"}}},
		{input: "1..1",
			want: []Token{{Type: Number, Value: "1."}, {Type: Number, Value: ".1"}}},
		{input: "0x",
			want: []Token{{Type: Error, Value: "illegal number"}}},
		{input: "0xHHHH",
			want: []Token{{Type: Error, Value: "illegal number"}}},
		{input: "1e1",
			want: []Token{{Type: Number, Value: "1e1"}}},
		{input: "1.1e1",
			want: []Token{{Type: Number, Value: "1.1e1"}}},
		{input: "32423432.234234234e1",
			want: []Token{{Type: Number, Value: "32423432.234234234e1"}}},
		{input: "1.1e1.123",
			want: []Token{{Type: Number, Value: "1.1e1.123"}}},
		{input: ".2e.3",
			want: []Token{{Type: Number, Value: ".2e.3"}}},
		{input: "0.2e0.3",
			want: []Token{{Type: Number, Value: "0.2e0.3"}}},
		{input: "-0.2e+0.3",
			want: []Token{{Type: Number, Value: "-0.2e+0.3"}}},
		{input: "+0.2e-0.3",
			want: []Token{{Type: Number, Value: "+0.2e-0.3"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			tokens := Tokenize(testCase.input)
			if !reflect.DeepEqual(tokens, testCase.want) {
				t.Fatalf("want tokens %+v, got %+v", testCase.want, tokens)
			}
		})
	}
}
