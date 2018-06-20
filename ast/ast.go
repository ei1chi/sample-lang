package ast

import (
	"strings"

	"github.com/ei1chi/sample-lang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}

	return out.String()
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

func (l *LetStmt) String() string {
	var out strings.Builder

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

func (l *LetStmt) stmtNode() {}

func (l *LetStmt) TokenLiteral() string {
	return l.Token.Literal
}

type ReturnStmt struct {
	Token       token.Token
	ReturnValue Expr
}

func (r *ReturnStmt) String() string {
	var out strings.Builder

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

func (r *ReturnStmt) stmtNode() {}

func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Literal
}

type ExprStmt struct {
	Token token.Token // 式の最初のトークン
	Expr  Expr
}

func (e *ExprStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

func (e *ExprStmt) stmtNode() {}

func (e *ExprStmt) TokenLiteral() string {
	return e.Token.Literal
}

type Ident struct {
	Token token.Token // == token.IDENT
	Value string
}

func (i *Ident) String() string {
	return i.Value
}

func (i *Ident) exprNode() {}

func (i *Ident) TokenLiteral() string {
	return i.Token.Literal
}
