package lexer

import (
	"testing"

	"github.com/ei1chi/sample-lang/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5; !-/*; if a return true else false; 10 == 10; 10 != 9;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.ELSE, "else"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
	}

	l := NewLexer(input)
	for i, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, test.expectedType, tok.Type)
		}
		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, test.expectedLiteral, tok.Literal)
		}
	}
}
