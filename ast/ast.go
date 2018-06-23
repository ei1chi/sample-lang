package ast

import (
	"bytes"
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

type BlockStmt struct {
	Token token.Token // "{" token
	Stmts []Stmt
}

func (b *BlockStmt) stmtNode() {}

func (b *BlockStmt) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStmt) String() string {
	var out strings.Builder

	for _, stmt := range b.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()
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

// =========================================
// Expressions
// =========================================

type PrefixExpr struct {
	Token    token.Token
	Operator string
	Right    Expr
}

func (p *PrefixExpr) exprNode() {}

func (p *PrefixExpr) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpr struct {
	Token    token.Token
	Left     Expr
	Operator string
	Right    Expr
}

func (i *InfixExpr) exprNode() {}

func (i *InfixExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpr) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
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

type IntLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntLiteral) String() string {
	return i.Token.Literal
}

func (i *IntLiteral) exprNode() {}

func (i *IntLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) exprNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpr struct {
	Token token.Token // == "if"
	Cond  Expr
	Cons  *BlockStmt
	Alt   *BlockStmt
}

func (i *IfExpr) exprNode() {}

func (i *IfExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpr) String() string {
	var out strings.Builder

	out.WriteString("if")
	out.WriteString(i.Cond.String())
	out.WriteString(" ")
	out.WriteString(i.Cons.String())

	if i.Alt != nil {
		out.WriteString("else ")
		out.WriteString(i.Alt.String())
	}

	return out.String()
}
