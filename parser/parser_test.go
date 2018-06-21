package parser

import (
	"fmt"
	"testing"

	"github.com/ei1chi/sample-lang/ast"
	"github.com/ei1chi/sample-lang/lexer"
)

func TestLetStmts(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Stmts) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Stmts))
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		stmt := program.Stmts[i]
		if !testLetStmt(t, stmt, test.expectedIdent) {
			return
		}
	}
}

func testLetStmt(t *testing.T, s ast.Stmt, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s not *ast.LetStmt. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStmts(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 3 {
		t.Fatalf("progarm.Stmts does not contain 3 statements. got=%d", len(program.Stmts))
	}

	for _, stmt := range program.Stmts {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt not *ast.returnStmt. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentExpr(t *testing.T) {
	input := `foobar;`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T", program.Stmts[0])
	}

	ident, ok := stmt.Expr.(*ast.Ident)
	if !ok {
		t.Fatalf("exp not *ast.Ident. got=%T", stmt.Expr)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntLiteralExpr(t *testing.T) {
	input := `5;`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T", program.Stmts[0])
	}

	literal, ok := stmt.Expr.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntLiteral. got=%T", stmt.Expr)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %s. got=%d", "5", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpr(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, test := range prefixTests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Stmts))
		}
		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T", program.Stmts[0])
		}

		exp, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpr. got=%T", stmt.Expr)
		}
		if exp.Operator != test.operator {
			t.Fatalf("exp.Operator is not '%s' got=%s", test.operator, exp.Operator)
		}
		if !testIntLiteral(t, exp.Right, test.intValue) {
			return
		}
	}
}

func testIntLiteral(t *testing.T, il ast.Expr, value int64) bool {
	integ, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("il not *ast.IntLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpr(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("program.Stmt[0] is not ast.ExprStmt. got=%T", program.Stmts[0])
		}

		exp, ok := stmt.Expr.(*ast.InfixExpr)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpr. got=%T", stmt.Expr)
		}

		if !testIntLiteral(t, exp.Left, test.leftValue) {
			return
		}

		if exp.Operator != test.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", test.operator, exp.Operator)
		}

		if !testIntLiteral(t, exp.Right, test.rightValue) {
			return
		}
	}

}
