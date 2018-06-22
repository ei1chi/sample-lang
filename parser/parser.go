package parser

import (
	"fmt"
	"strconv"

	"github.com/ei1chi/sample-lang/ast"
	"github.com/ei1chi/sample-lang/lexer"
	"github.com/ei1chi/sample-lang/token"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// ２つトークンを読み込み、curTokenとpeekTokenの両方に値をセットする
	p.nextToken()
	p.nextToken()

	// register prefixes
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	prefixes := map[token.TokenType]prefixParseFn{
		token.IDENT: p.parseIdent,
		token.INT:   p.parseIntLiteral,
		token.BANG:  p.parsePrefixExpr,
		token.MINUS: p.parsePrefixExpr,
	}
	for tok, fn := range prefixes {
		p.prefixParseFns[tok] = fn
	}

	// register infixes
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	infixes := map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpr,
		token.MINUS:    p.parseInfixExpr,
		token.SLASH:    p.parseInfixExpr,
		token.ASTERISK: p.parseInfixExpr,
		token.EQ:       p.parseInfixExpr,
		token.NOT_EQ:   p.parseInfixExpr,
		token.LT:       p.parseInfixExpr,
		token.GT:       p.parseInfixExpr,
	}
	for tok, fn := range infixes {
		p.infixParseFns[tok] = fn
	}

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Stmts = []ast.Stmt{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.curToken}

	stmt.Expr = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// =========================================
// Expressions
// =========================================

const (
	_ int = iota
	LOWEST
	EQUALS
	COMPARE
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precs = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       COMPARE,
	token.GT:       COMPARE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (p *Parser) peekPrec() int {
	if p, ok := precs[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrec() int {
	if p, ok := precs[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// parseExpr lowerPrec より優先度の高い演算子があるかを調べて繋げる。
func (p *Parser) parseExpr(lowerPrec int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	exp := prefix()

	for {
		// セミコロンが来たら解析終了
		if p.peekTokenIs(token.SEMICOLON) {
			break
		}

		// 受け取ったprec以下なら解析中断
		// SUMの右側に続く式なら、PRODUCT以上を期待する
		// 次にSUMが来たら中断して処理を上流に返す
		if p.peekPrec() <= lowerPrec {
			break
		}

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			break
		}

		p.nextToken()

		exp = infix(exp)
	}

	return exp
}

func (p *Parser) parseIdent() ast.Expr {
	return &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntLiteral() ast.Expr {
	lit := &ast.IntLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expr.Right = p.parseExpr(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	prec := p.curPrec()
	p.nextToken()
	expr.Right = p.parseExpr(prec)

	return expr
}
