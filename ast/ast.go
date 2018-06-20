package ast

import "github.com/ei1chi/sample-lang/token"

type Node interface {
	TokenLiteral() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	}
	return ""
}

type LetStmt struct {
	Token token.Token
	Name  *Ident
	Value Expr
}

func (l *LetStmt) stmtNode() {}

func (l *LetStmt) TokenLiteral() string {
	return l.Token.Literal
}

type ReturnStmt struct {
	Token       token.Token
	ReturnValue Expr
}

func (r *ReturnStmt) stmtNode() {}

func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Literal
}

type Ident struct {
	Token token.Token // == token.IDENT
	Value string
}

func (i *Ident) exprNode() {}

func (i *Ident) TokenLiteral() string {
	return i.Token.Literal
}
