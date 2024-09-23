package lexer

import (
	"testing"

	"github.com/kijimaD/nov/token"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `*label
こんにちは[l]あああ
←改行した。[p]
[image source="test.png" page="fore"]
[wait time="100"]`
	l := NewLexer(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{
			expectedType:    token.ASTERISK,
			expectedLiteral: "*",
		},
		{
			expectedType:    token.TEXT,
			expectedLiteral: "label",
		},
		{
			expectedType:    token.NEWLINE,
			expectedLiteral: "\n",
		},
		{
			expectedType:    token.TEXT,
			expectedLiteral: "こんにちは",
		},
		{
			expectedType:    token.LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "l",
		},
		{
			expectedType:    token.RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    token.TEXT,
			expectedLiteral: "あああ",
		},
		{
			expectedType:    token.NEWLINE,
			expectedLiteral: "\n",
		},
		{
			expectedType:    token.TEXT,
			expectedLiteral: "←改行した。",
		},
		{
			expectedType:    token.LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "p",
		},
		{
			expectedType:    token.RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    token.NEWLINE,
			expectedLiteral: "\n",
		},
		{
			expectedType:    token.LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "image",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "source",
		},
		{
			expectedType:    token.EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    token.STRING,
			expectedLiteral: "test.png",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "page",
		},
		{
			expectedType:    token.EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    token.STRING,
			expectedLiteral: "fore",
		},
		{
			expectedType:    token.RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    token.NEWLINE,
			expectedLiteral: "\n",
		},
		{
			expectedType:    token.LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "wait",
		},
		{
			expectedType:    token.IDENT,
			expectedLiteral: "time",
		},
		{
			expectedType:    token.EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    token.STRING,
			expectedLiteral: "100",
		},
		{
			expectedType:    token.RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    token.EOF,
			expectedLiteral: "",
		},
	}

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}
