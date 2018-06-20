package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/ei1chi/sample-lang/token"
)

type Lexer struct {
	input   string
	pos     int  // 現在の位置
	readPos int  // 現在の文字の次
	ch      rune // 現在検査中の文字
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	r, size := utf8.DecodeRuneInString(l.input[l.readPos:])
	if size == 0 {
		l.ch = 0
	} else {
		l.ch = r
	}
	l.pos = l.readPos
	l.readPos += size
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdent() string {
	begin := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[begin:l.pos]
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	begin := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[begin:l.pos]
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
