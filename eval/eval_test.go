package eval

import (
	"testing"

	"github.com/ei1chi/sample-lang/lexer"
	"github.com/ei1chi/sample-lang/object"
	"github.com/ei1chi/sample-lang/parser"
)

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	return Eval(program)
}

func TestEvalBooleanExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true != false", true},
		{"true == false", false},
		{"(1 < 2) == true", true},
	}

	for _, test := range tests {
		evaled := testEval(test.input)
		testBooleanObject(t, evaled, test.expected)
	}
}

func TestIfElseExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, test := range tests {
		evaled := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaled, int64(integer))
		} else {
			testNullObject(t, evaled)
		}
	}
}

func TestReturnStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}

			return 1;
		}
		`, 10},
	}

	for _, test := range tests {
		evaled := testEval(test.input)
		testIntegerObject(t, evaled, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, test := range tests {
		evaled := testEval(test.input)

		errObj, ok := evaled.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaled, evaled)
			continue
		}

		if errObj.Message != test.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", test.expectedMessage, errObj.Message)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("obejct is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestEvalIntExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2*2*2", 8},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"3 * (3 * 3) + 10", 37},
	}

	for _, test := range tests {
		evaled := testEval(test.input)
		testIntegerObject(t, evaled, test.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaled := testEval(test.input)
		testBooleanObject(t, evaled, test.expected)
	}
}
